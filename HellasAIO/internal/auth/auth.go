package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/loading"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/internal/version"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
	"github.com/jaypipes/ghw"
	"github.com/valyala/fastjson"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

var (
	AuthToken string
)

func newSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateHWID() string {
	block, _ := ghw.Block()
	var disks []string
	for _, disk := range block.Disks {
		disks = append(disks, disk.SerialNumber)
	}

	userStruct, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username := userStruct.Username

	return newSHA256(strings.Join(disks, ",") + "," + username)
}

func validateKey(key string) bool {
	authClient, _ := hclient.NewClient()

	requestBody := AuthRequest{
		Key:  key,
		HWID: generateHWID(),
	}

	_, err := authClient.NewRequest().
		SetURL("https://api.hellasaio.com/api/auth").
		SetMethod("POST").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "*/*").
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		log.Fatalln("Could not request auth server, shutting down.")
	}

	if fastjson.GetBool(authClient.LatestResponse.Body(), "success") {
		AuthToken = fastjson.GetString(authClient.LatestResponse.Body(), "access_token")
		return true
	} else {
		log.Fatalln(fastjson.GetString(authClient.LatestResponse.Body(), "message"))
	}

	return false
}

func getVersion() string {
	authClient, _ := hclient.NewClient()

	_, err := authClient.NewRequest().
		SetURL("https://api.hellasaio.com/api/latest").
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "*/*").
		Do()

	if err != nil {
		log.Fatalln(err.Error())
	}

	return fastjson.GetString(authClient.LatestResponse.Body(), "version")
}

func Initialize() {
	var key string
	if len(loading.Data.Settings.Settings.AuthKey) != 30 {
		fmt.Println("Paste your auth key & press enter: ")
		fmt.Scanln(&key)
	} else {
		key = loading.Data.Settings.Settings.AuthKey
	}

	if validateKey(key) {
		if key == "devkey" {
			utils.Debug = func() bool {
				if len(os.Args) > 1 {
					if os.Args[1] == "--debug" {
						return true
					}
				}
				return false
			}()
		}

		loading.Data.Settings.Settings.AuthKey = key
		loading.Data.Settings.Settings.Save()
		fmt.Println("Welcome!")
		fmt.Println("HellasAIO " + version.Version)
		go func() {
			for {
				time.Sleep(time.Second * 30)
				validateKey(key)
			}
		}()
	} else {
		log.Fatal("Invalid auth key")
	}
}
