package cloudflare

import (
	"bytes"
	"github.com/HellasAIO/HellasAIO/internal/task"
)

func DetectCloudflare(task *task.Task) bool {
	/*
		if m.Client.LatestResponse.StatusCode() == 503 || bytes.Contains(m.Client.LatestResponse.Body(), []byte("cf-challenge-running")) {
				logs.Log(m, "Fuel is down, or cloudflare is blocking us.")
				logs.Log(m, "Solving cloudflare...")
				success := cloudflare.GetClearanceCookie(m, "https://www.fuel.com.gr", f.ProxyURL)
				if success {
					logs.Log(m, "Cloudflare solved.")
					return LOGIN
				} else {
					logs.Log(m, "Failed to solve cloudflare.")
					time.Sleep(m.Delay)
					return LOGIN
				}
			}
	*/

	responseCode := task.Client.LatestResponse.StatusCode()
	responseBody := task.Client.LatestResponse.Body()

	return responseCode == 403 && bytes.Contains(responseBody, []byte("1020")) || responseCode == 503 && bytes.Contains(responseBody, []byte("jschal_js")) || responseBody == nil
}
