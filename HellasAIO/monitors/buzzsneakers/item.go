package buzzsneakersmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/valyala/fastjson"
	"net/url"
	"time"
)

func getItem(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	requestBody := url.Values{}
	requestBody.Add("nbAjax", "1")
	requestBody.Add("task", "live_search")
	requestBody.Add("returnType", "json")
	requestBody.Add("query", m.Product)
	requestBody.Add("typeSearch", "product")

	_, err := m.Client.NewRequest().
		SetURL("https://www.buzzsneakers.gr/athlitika-papoutsia/").
		SetMethod("POST").
		SetDefaultHeadersBuzz().
		SetFormBody(requestBody).
		Do()

	if err != nil {
		logs.Log(m, "Error getting item.")
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	return handleItemRequest(m, b)
}

func handleItemRequest(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	if m.Client.LatestResponse.StatusCode() != 200 {
		logs.Log(m, "Error getting item.")
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	if len(fastjson.GetString(m.Client.LatestResponse.Body(), "info")) == 28 || len(fastjson.GetString(m.Client.LatestResponse.Body(), "info")) == 68 {
		logs.Log(m, "Could not find MSKU.")
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	b.ProductId = utils.GetBuzzProductId(m.Client.LatestResponse.BodyAsString())
	logs.Log(m, "Got product id: "+b.ProductId)
	return GET_PRODUCT_INFO
}
