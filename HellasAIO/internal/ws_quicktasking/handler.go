package ws_quicktasking

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/auth"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/avast/retry-go"
	"github.com/getsentry/sentry-go"
	"github.com/valyala/fastjson"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

func makeTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		VerifyConnection: func(state tls.ConnectionState) error {
			for _, peercert := range state.PeerCertificates {
				der, err := x509.MarshalPKIXPublicKey(peercert.PublicKey)
				if err != nil {
					log.Fatalln("Failed to get public key (ws).")
				}

				var DNSName string
				if len(peercert.DNSNames) > 0 {
					DNSName = peercert.DNSNames[0]
				} else {
					DNSName = "Unknown Site"
				}

				hash := sha256.Sum256(der)
				stringHash := fmt.Sprintf("%x", hash)

				if utils.Debug {
					fmt.Println(fmt.Sprintf("%s: %s", DNSName, stringHash))
				}

				if stringHash == "6b68e323e3afd92a81513d1d89528893dc822c1501890afc34607c7ed9d33032" {
					return nil
				} else {
					log.Fatalln("HellasAIO Auth Server: SSL Missmatch")
				}
			}
			return errors.New("invalid certificate")
		},
	}
}

func makeTransport() *http.Transport {
	return &http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig:   makeTLSConfig(),
	}
}

func handleWebsocket(success chan bool) {
	defer sentry.Recover()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var authed = false

	client := &http.Client{Timeout: 15 * time.Second}
	client.Transport = makeTransport()

	var c *websocket.Conn
	options := websocket.DialOptions{
		HTTPClient: client,
	}

	defer log.Fatalln("Tried to reconnect 10 times to websocket server, but failed. Closing bot...")

	_ = retry.Do(func() error {
		defer time.Sleep(1 * time.Second)
		c, _, err = websocket.Dial(ctx, "wss://api.hellasaio.com/api/ws?token="+auth.AuthToken, &options)
		if err != nil {
			fmt.Println("Failed to connect to websocket server. Retrying...")
			return err
		} else {
			fmt.Println("Successfully connected to quicktask websocket.")
		}

		for {
			_, message, err := c.Read(ctx)
			if err != nil {
				if errors.Is(err, websocket.CloseError{Code: websocket.StatusPolicyViolation, Reason: ""}) {
					log.Fatalln("Failed to authenticate to websocket server.")
				} else {
					log.Println("Error getting websocket message.")
					return err
				}
			}

			if authed == false {
				if fastjson.GetBool(message, "success") {
					go func() { success <- true }()
					authed = true
				}
			}

			if authed {
				if fastjson.Exists(message, "siteId") {
					go handleQuicktaskMessage(message)
				}
			}
		}
	}, retry.Attempts(10), retry.MaxDelay(15*time.Second), retry.RetryIf(func(err error) bool {
		return ctx.Err() == nil && !errors.Is(err, context.Canceled)
	}))
}
