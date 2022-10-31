package fuelmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/cloudflare"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"time"
)

func solveCloudflare(m *task.Task, f *FuelInternal) task.TaskState {
	didExist := m.Client.InitSessionJar(f.Account)

	success := cloudflare.GetClearanceCookie(m, "https://www.fuel.com.gr", f.ProxyURL)
	if success {
		logs.Log(m, "Cloudflare solved.")
	} else {
		logs.Log(m, "Failed to solve cloudflare.")
		time.Sleep(m.Delay)
		return SOLVE_CLOUDFLARE
	}

	if m.Mode == "login" {
		return LOGIN
	}

	if didExist {
		logs.Log(m, "Skipping login using saved session.")
		return GET_CART_ID
	}

	return LOGIN
}
