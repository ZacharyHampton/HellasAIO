package buzzsneakersmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
)

func initialize(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	if !utils.Contains([]string{"msku", "sku"}, m.Mode) {
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

	m.CheckoutData.Website = "buzzsneakers"
	m.CheckoutData.Mode = m.Mode
	m.CheckoutData.ProductMSKU = m.Product
	m.Client = client

	if m.Mode == "sku" {
		b.ProductId = m.Product
		return GET_PRODUCT_INFO
	}

	return GET_ITEM
}
