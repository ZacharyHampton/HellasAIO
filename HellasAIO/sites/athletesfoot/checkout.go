package athletesfoot

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"strconv"
	"time"
)

func checkout(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	postData := "orderId=" + strconv.Itoa(i.OrderId) + "&store=0&deliveryPoint=1&shippingMethod=%7B%22MethodId%22%3A1%2C%22MethodName%22%3A%22%CE%91%CF%80%CE%BF%CF%83%CF%84%CE%BF%CE%BB%CE%AE+%CE%BC%CE%AD%CF%83%CF%89+courrier%22%2C%22MethodDescription%22%3A%22%22%2C%22CompanyId%22%3A9%2C%22Duration%22%3A%224-10%22%2C%22AdditionalDuration%22%3A%22%22%2C%22DurationText%22%3A%224-10+%CE%B5%CF%81%CE%B3%CE%AC%CF%83%CE%B9%CE%BC%CE%B5%CF%82+%CE%BC%CE%AD%CF%81%CE%B5%CF%82%22%2C%22AdditionalDurationText%22%3A%22%22%2C%22CompanyName%22%3A%22Speedex%22%2C%22CompanyPhones%22%3A%22%22%2C%22CompanyAddress%22%3A%22%22%2C%22UserCompany%22%3A%7B%22Name%22%3Anull%2C%22Address%22%3Anull%2C%22Town%22%3Anull%2C%22Area%22%3Anull%2C%22Phone%22%3Anull%2C%22PostalCode%22%3Anull%7D%2C%22PartnerDurationText%22%3A%22%22%2C%22FromStoreTimeText%22%3A%22%22%2C%22PaymentMethods%22%3A%5B3%2C1%2C5%5D%2C%22AvailableShippingCompanies%22%3A%5B17%2C18%2C9%5D%7D&paymentMethod=%7B%22Name%22%3A%22%CE%91%CE%BD%CF%84%CE%B9%CE%BA%CE%B1%CF%84%CE%B1%CE%B2%CE%BF%CE%BB%CE%AE%22%2C%22Code%22%3A%22cashOnDelivery%22%2C%22MethodId%22%3A3%2C%22ProcessorId%22%3A5%2C%22InstallmentsCount%22%3A0%7D&orderType=%7B%22Id%22%3A1%2C%22IsVisible%22%3Atrue%2C%22Code%22%3A%22Normal%22%2C%22Description%22%3A%22%CE%9A%CE%B1%CE%BD%CE%BF%CE%BD%CE%B9%CE%BA%CE%AE%22%2C%22DescriptionEl%22%3A%22%CE%9A%CE%B1%CE%BD%CE%BF%CE%BD%CE%B9%CE%BA%CE%AE%22%2C%22DescriptionEn%22%3A%22Normal%22%2C%22IsPaymentOrderType%22%3Afalse%7D"

	_, err := t.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/engine/processorder").
		SetMethod("POST").
		SetHeader("Host", "www.theathletesfoot.gr").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Sec-Ch-Ua", `"Chromium";v="93", " Not;A Brand";v="99"`).
		SetHeader("Sec-Ch-Ua-Mobile", `?0`).
		SetHeader("Sec-Ch-Ua-Platform", `"Windows"`).
		SetHeader("Upgrade-Insecure-Requests", `1`).
		SetHeader("Origin", `https://www.theathletesfoot.gr`).
		SetHeader("Content-Type", `application/x-www-form-urlencoded`).
		SetHeader("User-Agent", userAgent).
		SetHeader("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`).
		SetHeader("Sec-Fetch-Site", "same-origin").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Referer", "https://www.theathletesfoot.gr/pages/checkout/default.aspx?lang=el").
		SetHeader("Accept-Encoding", "gzip, deflate").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Connection", "close").
		SetBody(postData).
		Do()

	if err != nil {
		logs.Log(t, "Failed to send checkout order request.")
		time.Sleep(t.Delay)
		return CHECKOUT_ORDER
	}

	return handleCheckoutResponse(t, i)
}

func handleCheckoutResponse(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	if t.Client.LatestResponse.StatusCode() != 200 {
		logs.Log(t, "Failed to checkout order.")
		title.AddFailure()
		time.Sleep(t.Delay)
		return CHECKOUT_ORDER
	}

	logs.Log(t, "Checked out order successfully.")
	title.AddCheckout()
	t.CheckoutData.Status = "success"
	return task.DoneTaskState
}
