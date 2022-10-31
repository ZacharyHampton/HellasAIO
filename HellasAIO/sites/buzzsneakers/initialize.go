package buzzsneakers

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
)

func initialize(m *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	if !utils.Contains([]string{"login", "normal"}, m.Mode) {
		logs.Log(m, "Mode is not supported for this site. (invalid mode)")
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
		logs.Log(m, "Failed to create client. (error: %s)", err.Error())
		return task.ErrorTaskState
	}

	if m.AccountId != "" {
		b.Account, err = account.GetAccount(m.SiteId, m.AccountId)
		if err != nil {
			logs.Log(m, "Failed to get account. (error: %s)", err.Error())
			return task.ErrorTaskState
		}
	} else {
		logs.Log(m, "no account specified")
		return task.ErrorTaskState
	}

	didExist := client.InitSessionJar(b.Account)
	m.Client = client
	if m.Mode == "login" {
		return LOGIN
	}

	if didExist {
		logs.Log(m, "Skipping login using saved session.")
		return WAIT_FOR_MONITOR
	}

	return LOGIN
}
