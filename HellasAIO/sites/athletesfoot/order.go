package athletesfoot

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"time"
)

func getOrderId(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	_, err := t.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/pages/checkout/default.aspx?lang=el").
		SetMethod("GET").
		SetHeader("user-agent", userAgent).
		SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Do()

	if err != nil {
		logs.Log(t, "Failed to get order ID")
		time.Sleep(t.Delay)
		return GET_ORDER_ID
	}

	return handleOrderIDResponse(t, i)
}

func handleOrderIDResponse(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	orderId := utils.GetAFOrderID(t.Client.LatestResponse.BodyAsString())
	if orderId == -1 {
		// handle error and retry
		logs.Log(t, "orderId == -1")
		time.Sleep(t.Delay)
		return GET_ORDER_ID
	}

	i.OrderId = orderId
	return SUBMIT_ORDER
}
