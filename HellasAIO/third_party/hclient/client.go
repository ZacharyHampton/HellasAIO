package hclient

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/sessions"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	sessionjar "github.com/juju/persistent-cookiejar"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	NoCookieJarErr = errors.New("no cookie jar in client")
)

// NewClient creates a new http client
// Takes in the optional arguments: proxy, servername
func NewClient(parameters ...string) (*Client, error) {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyConnection: func(state tls.ConnectionState) error {
			for _, peercert := range state.PeerCertificates {
				der, err := x509.MarshalPKIXPublicKey(peercert.PublicKey)
				if err != nil {
					log.Println("Failed to get public key (https).")
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

				if fingerprints[stringHash] == 1 {
					return nil
				} else {
					fmt.Println(DNSName + ": SSL mismatch.")
				}
			}
			return fmt.Errorf("invalid certificate")
		},
	}

	// parameters[0] = proxy
	// parameters[1] = sni
	if len(parameters) > 1 && len(parameters[1]) > 0 {
		tlsClientConfig.ServerName = parameters[1]
	}

	transport := &http.Transport{
		ForceAttemptHTTP2: true,
		TLSClientConfig:   tlsClientConfig,
	}

	if len(parameters) > 0 && len(parameters[0]) > 0 {
		proxyUrl, _ := url.Parse(parameters[0])

		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return &Client{
		client: &http.Client{
			Transport: transport,
		},
		LatestResponse: &Response{},
	}, nil
}

// NewRequest creates a new request under a specified http client
func (c *Client) NewRequest() *Request {
	return &Request{
		client: c,
		header: make(http.Header),
	}
}

func (c *Client) InitCookieJar() {
	if c.client.Jar == nil {
		c.client.Jar, _ = cookiejar.New(nil)
	}
}

// InitSessionJar creates session jar, returns if it already existed or not
func (c *Client) InitSessionJar(account *account.Account) bool {
	didExist := sessions.DoesSessionExist(account)

	jar, err := sessionjar.New(&sessionjar.Options{
		Filename: fmt.Sprintf("../.sessions/%s/%s.sessions", strings.Replace(utils.SiteIDtoSiteString[account.SiteId], "@", "", -1), account.Email),
	})

	if err != nil {
		fmt.Println("Failed to initialize session. ", err)
		return false
	}

	c.jar = jar
	c.client.Jar = jar
	return didExist
}

func (c *Client) SaveCookies() {
	if c.client.Jar != nil {
		err := c.jar.Save()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// AddCookie adds a new cookie to the request client cookie jar
func (c *Client) AddCookie(u *url.URL, cookie *http.Cookie) error {
	if c.client.Jar == nil {
		c.client.Jar, _ = cookiejar.New(nil)
	}

	currentCookies := c.client.Jar.Cookies(u)
	currentCookies = append(currentCookies, cookie)
	c.client.Jar.SetCookies(u, currentCookies)

	return nil
}

// RemoveCookie removes the specified cookie from the request client cookie jar
func (c *Client) RemoveCookie(u *url.URL, cookie string) error {
	if c.client.Jar == nil {
		c.client.Jar, _ = cookiejar.New(nil)
	}

	newCookie := &http.Cookie{
		Name:  cookie,
		Value: "",
	}

	c.client.Jar.SetCookies(u, []*http.Cookie{newCookie})

	return nil
}

func (c *Client) AddCookieByName(r *Response, u *url.URL, name string) error {
	cookie := r.GetCookieByName(name)
	if cookie != nil {
		err := c.AddCookie(u, cookie)
		if err != nil {
			return err
		}
	}

	return nil
}

// Do will send the specified request
func (c *Client) Do(r *http.Request) (*Response, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// https://help.socketlabs.com/docs/how-to-fix-error-only-one-usage-of-each-socket-address-protocolnetwork-addressport-is-normally-permitted
	// https://www.geeksforgeeks.org/http-headers-connection/#:~:text=close%20This%20close%20connection%20directive,want%20your%20connection%20to%20close.
	r.Close = true // perhaps set this to false?

	response := &Response{
		headers:    resp.Header,
		body:       body,
		status:     resp.Status,
		statusCode: resp.StatusCode,
		cookies:    resp.Cookies(),
	}

	c.LatestResponse = response
	if utils.Debug {
		fmt.Println(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
		fmt.Println(fmt.Sprintf("Response Body: %s", response.BodyAsString()))
	}

	return response, nil
}
