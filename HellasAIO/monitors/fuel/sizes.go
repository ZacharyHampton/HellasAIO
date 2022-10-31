package fuelmonitor

import (
	"bytes"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/cloudflare"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/valyala/fastjson"
	"math/rand"
	"strings"
	"time"
)

func findSizes(m *task.Task, f *FuelInternal) task.TaskState {
	requestBody := SizesRequest{
		Query:     "query getProducts($msku: String!){\r\n  products(filter: {sku: {eq: $msku}}) {\r\n    items {\r\n      name\r\n      sku\r\n      id\r\n      stock_status\r\n      __typename\r\n      ... on ConfigurableProduct {\r\n        variants {\r\n            product {\r\n                id\r\n                created_at\r\n                websites {\r\n                    name\r\n                    id\r\n                }\r\n                __typename\r\n                id\r\n                upcoming \r\n                name\r\n                sku\r\n                stock_status\r\n                is_raffle_item\r\n                only_x_left_in_stock\r\n                size\r\n            }\r\n        }\r\n      }\r\n    }\r\n  }\r\n}\r\n",
		Variables: SizesRequestVariables{MSKU: f.ParentSKU},
	}

	_, err := m.Client.NewRequest().
		SetURL("https://www.fuel.com.gr/el/graphql").
		SetMethod("POST").
		SetHeader("User-Agent", utils.UserAgent).
		SetHeader("Content-Type", "application/json").
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		return GET_SIZE
	}

	return handleSizeResponse(m, f)
}

func handleSizeResponse(m *task.Task, f *FuelInternal) task.TaskState {
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
		logs.Log(m, "Unknown response.")
		time.Sleep(m.Delay)
		return GET_SIZE
	}

	if !(bytes.Contains(m.Client.LatestResponse.Body(), []byte(m.Product))) {
		logs.Log(m, "Product not found")
		time.Sleep(m.Delay)
		return GET_SIZE
	}

	jValue, err := fastjson.ParseBytes(m.Client.LatestResponse.Body())
	if err != nil {
		return GET_SIZE
	}

	var instockVariants []*fastjson.Value

	for _, value := range jValue.GetArray("data", "products", "items", "0", "variants") {
		if bytes.Equal(value.GetStringBytes("product", "stock_status"), []byte("IN_STOCK")) {
			instockVariants = append(instockVariants, value)
		}

	}

	// if there are instock variants
	if len(instockVariants) != 0 {
		rand.Seed(time.Now().UnixNano())
		variant := instockVariants[rand.Intn(len(instockVariants))]
		logs.Log(m, fmt.Sprintf("Found variant: %q", variant.GetStringBytes("product", "name")))
		f.VariantSKU = string(variant.GetStringBytes("product", "sku"))
		m.CheckoutData.Size = strings.Split(f.VariantSKU, " ")[1]
		f.ProductFound = true
	} else {
		logs.Log(m, "No instock variants found")
		return GET_SIZE
	}

	return CHECKOUT
}
