package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
)

func initialize(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if !utils.Contains([]string{"msku", "login"}, m.Mode) {
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
	client, err := hclient.NewClient(proxyURL)

	if err != nil {
		return task.ErrorTaskState
	}

	m.CheckoutData.Website = "athletesfoot"
	m.CheckoutData.Mode = m.Mode
	m.CheckoutData.ProductMSKU = m.Product

	if m.AccountId != "" {
		a.Account, _ = account.GetAccount(m.SiteId, m.AccountId)
	} else {
		logs.Log(m, "no account specified")
		return task.ErrorTaskState
	}

	m.Client = client

	didExist := client.InitSessionJar(a.Account)
	if m.Mode == "login" {
		return GET_SESSION
	}

	if didExist {
		logs.Log(m, "Skipping login using saved session.")
		return CLEAR_CART
	}

	return GET_SESSION
}
