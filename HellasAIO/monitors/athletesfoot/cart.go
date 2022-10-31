package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"time"
)

func clearCart(m *task.Task, a *AthletesFootInternal) task.TaskState {
	requestBody := GetCartRequest{
		CurrentItemId: 2,
		TemplateCode:  "homepage",
	}

	_, err := m.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/services/api/Cart/GetCartViewModel?lang=el").
		SetMethod("POST").
		SetDefaultHeadersAF().
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		logs.Log(m, "Error while making clearing cart request")
		time.Sleep(m.Delay)
		return CLEAR_CART
	}

	return HandleClearCartResponse(m, a)
}

func HandleClearCartResponse(m *task.Task, a *AthletesFootInternal) task.TaskState {
	// if cartclear successful, start checking stock, or else error or try again
	// we gonna do it blind, pray for no errors

	jValue, err := fastjson.ParseBytes(m.Client.LatestResponse.Body())
	if err != nil {
		logs.Log(m, "Error while parsing clear cart response")
		return CLEAR_CART
	}
	for _, value := range jValue.GetArray("data", "Cart", "Items") {
		requestBody2 := RemoveFromCartRequest{
			Id: value.GetInt("Id"),
		}

		_, err := m.Client.NewRequest().
			SetURL("https://www.theathletesfoot.gr/services/EcomService.svc/RemoveFromCart?lang=el").
			SetMethod("POST").
			SetDefaultHeadersAF().
			SetJSONBody(requestBody2).
			Do()

		if err != nil {
			logs.Log(m, "Error while clearing item.")
			time.Sleep(m.Delay)
			return CLEAR_CART
		}
	}

	return ADD_TO_CART
}
