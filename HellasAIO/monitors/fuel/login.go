package fuelmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/cloudflare"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"net/url"
	"time"
)

func login(m *task.Task, f *FuelInternal) task.TaskState {
	params := url.Values{}
	params.Add("form_key", "Rvq5JgU8qzsZaFNR")
	params.Add("login[username]", f.Account.Email)
	params.Add("login[password]", f.Account.Password)
	params.Add("send", "")

	_, err := m.Client.NewRequest().
		SetURL("https://www.fuel.com.gr/el/customer/account/loginPost/").
		SetMethod("POST").
		SetDefaultHeadersFuel().
		SetHost("www.fuel.com.gr").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Cookie", "cp_total_cart_items=0; cp_total_cart_value=0; cpab=9cf91df1-49dd-4e08-9156-09651da1f0a4; _gcl_au=1.1.1460003395.1655054482; mage-cache-storage=%7B%7D; mage-cache-storage-section-invalidation=%7B%7D; mage-translation-storage=%7B%7D; mage-translation-file-version=%7B%7D; _fbp=fb.2.1655054486365.347843962; cp_sessionTime=1655054475909; cto_bundle=NPWb3F9HVUluVjZWWHFMTnlxOTlkN1lSWGJURWk0UUp5OFM0bXNMWVQzakFaS2hDMCUyRm10NndENGlHN0FYVWRjbklxazFDbTE5SU8yS01oeXFzU09rSUxZUHB3UHNyUm16M3VFMTZBazJnQVMwckdlYmJEOERrQWJYSFpQQWtsRW5xR0NrTHJKR0l3UzlrU042RHBBWDFmcEElMkZnJTNEJTNE; PHPSESSID=ohfus5ios4bmpgmom3prc30osr; form_key=Rvq5JgU8qzsZaFNR; mage-messages=; recently_viewed_product=%7B%7D; recently_viewed_product_previous=%7B%7D; recently_compared_product=%7B%7D; recently_compared_product_previous=%7B%7D; product_data_storage=%7B%7D; mage-cache-sessid=true; user_allowed_save_cookie=%7B%221%22%3A1%7D; section_data_ids=%7B%22gtm%22%3A1655055736%2C%22cart%22%3Anull%2C%22customer%22%3Anull%2C%22messages%22%3Anull%2C%22captcha%22%3Anull%2C%22compare-products%22%3Anull%2C%22product_data_storage%22%3Anull%2C%22custom-minicart%22%3Anull%7D").
		// SetBody("form_key=Rvq5JgU8qzsZaFNR&login%5Busername%5D=" + f.Account.Email + "&login%5Bpassword%5D=" + f.Account.Password + "&send="). url encode password & email
		SetFormBody(params).
		Do()

	if err != nil {
		return LOGIN
	}

	return handleLoginResponse(m, f)
}

func handleLoginResponse(m *task.Task, f *FuelInternal) task.TaskState {
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

	if m.Client.LatestResponse.StatusCode() == 200 {
		formKey := m.Client.LatestResponse.GetCookieByName("form_key")
		mageMessages := m.Client.LatestResponse.GetCookieByName("mage-messages")

		if formKey == nil && mageMessages == nil {
			logs.Log(m, "Logged in!")
			m.Client.SaveCookies()
			if m.Mode == "login" {
				return task.DoneTaskState
			}
			return GET_CART_ID
		} else {
			mageMessagesString, _ := url.QueryUnescape(mageMessages.Value)
			logs.Log(m, "Failed to log in. Message: ", fastjson.GetString([]byte(mageMessagesString), "0", "text"))
			return task.ErrorTaskState
		}
	} else {
		logs.Log(m, "Error: ", m.Client.LatestResponse.StatusCode())
		time.Sleep(m.Delay)
	}

	return LOGIN
}
