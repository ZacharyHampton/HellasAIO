package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"github.com/valyala/fastjson"
	"log"
	"strings"
	"time"
)

func addToCart(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if strings.ToLower(m.Mode) == "msku" && a.foundItem == false {
		return GET_ITEM
	}

	requestBody := AddToCartRequest{}
	requestBody.Item.SKU = a.sku
	requestBody.Item.SubSKU = a.subSku
	requestBody.Item.PathId = a.pathId
	requestBody.Item.Quantity = "1.00" // make sure to change this to a changeable quantity
	requestBody.Item.ExtraAttributes = ""
	requestBody.Item.BundleItems = []string{}
	requestBody.Item.RecipeItems = []string{}
	requestBody.Item.ComboItems = []string{}
	requestBody.Item.IsDefaultRecipeQuantitiesIncluded = true
	requestBody.Item.EnhancedInfo.MinQuantity = 1
	requestBody.Item.EnhancedInfo.MaxQuantity = 1 // changeable
	requestBody.Item.EnhancedInfo.Quanta = 1      // changeable

	_, err := m.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/services/EcomService.svc/AddToCart?lang=el").
		SetMethod("POST").
		SetDefaultHeadersAF().
		SetHeader("Referer", "https://www.theathletesfoot.gr/andrika/papoutsia/ola-ta-papoutsia/nike-2-nike_718058/").
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		logs.Log(m, "Error while making add to cart request")
		return ADD_TO_CART
	}

	return handleATCResponse(m, a)
}

func handleATCResponse(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if fastjson.GetString(m.Client.LatestResponse.Body(), "d", "code") == "ValidationError" {
		// oos
		if m.Size != "random" {
			logs.Log(m, "OOS, finding new item.")
		} else {
			logs.Log(m, "OOS")
		}

		time.Sleep(m.Delay)
		a.foundItem = false
		return GET_ITEM
	}

	if !(fastjson.GetString(m.Client.LatestResponse.Body(), "d", "code") == "OperationSuccesful") {
		// oos or error
		logs.Log(m, "OOS or error while adding to cart")
		time.Sleep(m.Delay)
		return ADD_TO_CART
	}

	logs.Log(m, "Added to cart.")
	title.AddCart()
	logs.Log(m, "Notifying checkout tasks...")
	err := m.NotifyTasks(&AthletesFootMonitorData{Client: m.Client, CheckoutData: m.CheckoutData})
	if err != nil {
		log.Println("Failed to send notifications.")
	}

	logs.Log(m, "Checkout tasks notified. Sleeping for 60 seconds...")
	a.foundItem = false // reset found item
	time.Sleep(60 * time.Second)
	return CLEAR_CART
}
