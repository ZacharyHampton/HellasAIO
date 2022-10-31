package athletesfootmonitor

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"math/rand"
	"time"
)

func getItem(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if a.foundItem {
		return ADD_TO_CART
	}

	requestBody := GetItemRequest{
		EsFields:    "",
		EsGeoWindow: "",
		PageFrom:    1,
		PageSize:    1,
		Path:        "?hidefromsearch=1",
		Query:       m.Product, // msku
		Type:        "products",
	}

	_, err := m.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/engine/esqueryindex?lang=el").
		SetMethod("POST").
		SetDefaultHeadersAF().
		SetHeader("Referer", "https://www.theathletesfoot.gr/andrika/papoutsia/ola-ta-papoutsia/nike-2-nike_718058/").
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		logs.Log(m, "Error getting item", err)
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	return HandleItemResponse(m, a)
}

func HandleItemResponse(m *task.Task, a *AthletesFootInternal) task.TaskState {
	// validate product is found
	if !(fastjson.GetString(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Properties", "ManufacturerProductCode") == m.Product) {
		fmt.Println("Product not found")
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	a.foundItem = true
	a.sku = fastjson.GetInt(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Id")
	a.pathId = fastjson.GetInt(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Path")
	m.CheckoutData.ImageUrl = "https:" + fastjson.GetString(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Images", "0", "Src")
	m.CheckoutData.ProductName = fastjson.GetString(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Title")
	m.CheckoutData.Price = fastjson.GetFloat64(m.Client.LatestResponse.Body(), "hits", "hits", "0", "_source", "Price")
	jValue, err := fastjson.ParseBytes(m.Client.LatestResponse.Body())
	if err != nil {
		return GET_ITEM
	}

	var instockVariants []*fastjson.Value
	var sizeObjectsMap = make(map[string]*fastjson.Value)

	for _, value := range jValue.GetArray("hits", "hits", "0", "_source", "Dimensions") {
		if value.Exists("Stock") {
			instockVariants = append(instockVariants, value)
		}

		for _, dDList := range value.GetArray("DimensionDataList") {
			if string(dDList.GetStringBytes("DimensionType")) == "size" {
				sizeObjectsMap[string(dDList.GetStringBytes("Value"))] = value
			}
		}
	}

	switch m.Size {
	case "random":
		// if there are instock variants
		if len(instockVariants) != 0 {
			rand.Seed(time.Now().UnixNano())
			selectedValue := instockVariants[rand.Intn(len(instockVariants))]
			a.subSku = selectedValue.GetInt("SkuId")
			for _, dDList := range selectedValue.GetArray("DimensionDataList") {
				if string(dDList.GetStringBytes("DimensionType")) == "size" {
					m.CheckoutData.Size = string(dDList.GetStringBytes("Value"))
				}
			}
		} else if len(sizeObjectsMap) != 0 {
			// no instock variants, but there are sizes available
			// a.subSku = sizeObjectsMap[m.Size].GetInt("SkuId")
			return GET_ITEM
		} else {
			// no variants at all
			return GET_ITEM
		}
	default:
		if sizeObjectsMap[m.Size] != nil {
			a.subSku = sizeObjectsMap[m.Size].GetInt("SkuId")
		} else {
			return GET_ITEM
		}
	}

	return ADD_TO_CART
}
