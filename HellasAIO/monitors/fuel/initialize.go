package fuelmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
	"log"
)

func initialize(m *task.Task, f *FuelInternal) task.TaskState {
	if !utils.Contains([]string{"fast", "login"}, m.Mode) {
		logs.Log(m, "Mode is not supported for this site.")
		return task.ErrorTaskState
	}

	var proxyURL string
	if m.ProxyListID != "" {
		proxyObject, err := proxy.GetProxyFromProxyGroup(m.ProxyListID)
		if err == nil {
			proxyURL = proxyObject.URL
		} else {
			proxyURL = ""
		}
	}

	if proxyURL == "" {
		logs.Log(m, "Proxies are required on fuel. Please run with proxies (or enough proxies).")
		return task.ErrorTaskState
	}

	f.ProxyURL = proxyURL
	client, err := hclient.NewClient(proxyURL)

	if err != nil {
		return task.ErrorTaskState
	}

	if m.AccountId != "" {
		f.Account, err = account.GetAccount(m.SiteId, m.AccountId)
		if err != nil {
			log.Println("Error getting account: ", err.Error())
			return task.ErrorTaskState
		}
	} else {
		logs.Log(m, "no account specified")
		return task.ErrorTaskState
	}

	m.CheckoutData.Website = "fuel"
	m.CheckoutData.Mode = m.Mode
	m.CheckoutData.ProductMSKU = m.Product
	m.Client = client

	return SOLVE_CLOUDFLARE
}
