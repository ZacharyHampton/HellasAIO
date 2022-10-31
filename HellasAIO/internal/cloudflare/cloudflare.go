package cloudflare

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/auth"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
	"github.com/valyala/fastjson"
	"net/http"
	"net/url"
)

func GetClearanceCookie(t *task.Task, sUrl, proxyURL string) bool {
	fURL, err := url.Parse(sUrl)
	if err != nil {
		return false
	}

	if !DetectCloudflare(t) {
		logs.Log(t, "Cannot solve cloudflare: response was not from CF.")
		return false
	}

	if proxyURL == "" {
		logs.Log(t, "Cannot solve cloudflare: no proxies available.")
		return false
	}

	client, _ := hclient.NewClient()
	_, err = client.NewRequest().
		SetURL(fmt.Sprintf("https://api.hellasaio.com/api/cloudflare/%d", t.SiteId)).
		SetMethod("GET").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "*/*").
		SetHeader("User-Agent", utils.UserAgent).
		SetHeader("x-proxy", proxyURL).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", auth.AuthToken)).
		Do()

	if err != nil {
		return false
	}

	if !(fastjson.GetBool(client.LatestResponse.Body(), "success")) {
		logs.Log(t, "Cannot solve cloudflare: ", fastjson.GetString(client.LatestResponse.Body(), "message"))
		return false
	}

	responseJSON, _ := fastjson.Parse(client.LatestResponse.BodyAsString())
	cookies := responseJSON.GetObject("cookies")
	cookie := cookies.Get("cf_clearance")
	cookieBytes, err := cookie.StringBytes()
	if err != nil {
		return false
	}

	err = t.Client.AddCookie(fURL, &http.Cookie{
		Name:     "cf_clearance",
		Value:    string(cookieBytes),
		Path:     "/",
		Domain:   fURL.Host, // make website differential
		Secure:   true,
		HttpOnly: true,
		SameSite: 0,
	})
	if err != nil {
		return false
	}

	return true
}
