package fuelmonitor

import (
	"bytes"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/cloudflare"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/valyala/fastjson"
	"time"
)

func getCartId(m *task.Task, f *FuelInternal) task.TaskState {
	_, err := m.Client.NewRequest().
		SetURL("https://www.fuel.com.gr/el/graphql").
		SetMethod("POST").
		SetHeader("User-Agent", utils.UserAgent).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"query":"mutation {createEmptyCart}"}`).
		Do()

	if err != nil {
		logs.Log(m, fmt.Sprintf("Error getting cart id: %s", err.Error()))
		return GET_CART_ID
	}

	return handleCartIdResponse(m, f)
}

func handleCartIdResponse(m *task.Task, f *FuelInternal) task.TaskState {
	if cloudflare.DetectCloudflare(m) {
		logs.Log(m, "Cloudflare detected.")
		for {
			logs.Log(m, "Solving cloudflare...")
			success := cloudflare.GetClearanceCookie(m, "https://www.fuel.com.gr", f.ProxyURL)
			if success {
				logs.Log(m, "Cloudflare solved.")
				return GET_CART_ID
			} else {
				logs.Log(m, "Failed to solve cloudflare. Retrying...")
				time.Sleep(m.Delay)
			}
		}
	}

	if m.Client.LatestResponse.StatusCode() != 200 {
		logs.Log(m, fmt.Sprintf("Unknown response (%s). ", m.Client.LatestResponse.Status()))
		time.Sleep(m.Delay)
		return GET_CART_ID
	}

	if bytes.Contains(m.Client.LatestResponse.Body(), []byte("errors")) {
		logs.Log(m, "Error clearing cart.")
		time.Sleep(m.Delay)
		return GET_CART_ID
	}

	f.CartId = fastjson.GetString(m.Client.LatestResponse.Body(), "data", "createEmptyCart")
	f.ParentSKU = m.Product
	if m.Size == "random" {
		return GET_SIZE
	}

	// if not random
	m.CheckoutData.Size = m.Size
	f.VariantSKU = ParentToVariantConversion(f.ParentSKU, m.Size)
	f.ProductFound = true
	return CHECKOUT
}
