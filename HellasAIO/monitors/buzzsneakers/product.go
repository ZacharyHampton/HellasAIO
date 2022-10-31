package buzzsneakersmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"net/url"
	"strconv"
	"time"
)

func getProductInfo(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	logs.Log(m, "Getting product info...")

	requestBody := url.Values{}
	requestBody.Add("task", "getproductdata")
	requestBody.Add("productId", b.ProductId)
	requestBody.Add("nbAjax", "1")

	_, err := m.Client.NewRequest().
		SetURL("https://www.buzzsneakers.gr/athlitika-papoutsia/").
		SetMethod("POST").
		SetDefaultHeadersBuzz().
		SetFormBody(requestBody).
		Do()

	if err != nil {
		logs.Log(m, "Error getting product info.")
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	return handleProductInfoRequest(m, b)
}

func handleProductInfoRequest(m *task.Task, b *BuzzSneakersInternal) task.TaskState {
	if m.Client.LatestResponse.StatusCode() != 200 {
		logs.Log(m, "Error getting product info.")
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	if fastjson.GetString(m.Client.LatestResponse.Body(), "product", "productCode") != m.Product && m.Mode != "sku" {
		logs.Log(m, "Found invalid product.")
		time.Sleep(m.Delay)
		return GET_ITEM
	}

	products := BuzzProduct{
		InstockProducts: make([]BuzzCombinationProduct, 0),
		Products:        make([]BuzzCombinationProduct, 0),
	}

	jsonBody, err := fastjson.ParseBytes(m.Client.LatestResponse.Body())
	if err != nil {
		logs.Log(m, "Error parsing product info.")
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	m.CheckoutData.ImageUrl = "https://buzzsneakers.gr" + string(jsonBody.GetStringBytes("product", "image"))
	m.CheckoutData.ProductName = string(jsonBody.GetStringBytes("product", "name"))
	m.CheckoutData.ProductMSKU = fastjson.GetString(m.Client.LatestResponse.Body(), "product", "productCode")

	totalQuantity, err := strconv.ParseFloat(fastjson.GetString(m.Client.LatestResponse.Body(), "product", "quantity"), 64)
	if err != nil {
		logs.Log(m, "Error parsing product quantity.")
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	if totalQuantity <= 0 {
		logs.Log(m, "Product is out of stock.")
		time.Sleep(m.Delay)

		/*if m.Mode == "sku" {
			return GET_PRODUCT_INFO
		}*/

		// return GET_ITEM

		return GET_PRODUCT_INFO
	}

	var loadedSizeQuantity float64
	loadedSizeQuantity = 0

	for _, product := range jsonBody.GetArray("sizes") {
		quantity, err := strconv.ParseFloat(string(product.GetStringBytes("quantity")), 64)
		loadedSizeQuantity += quantity
		if err != nil {
			logs.Log(m, "Error parsing product quantity.")
			time.Sleep(m.Delay)
			return GET_PRODUCT_INFO
		}

		price, err := strconv.ParseFloat(string(product.GetStringBytes("salePriceWithTax")), 64)
		if err != nil {
			logs.Log(m, "Error parsing product price.")
			time.Sleep(m.Delay)
			return GET_PRODUCT_INFO
		}

		buzzObject := BuzzCombinationProduct{
			Quantity:      quantity,
			CombinationId: string(product.GetStringBytes("productCombinationId")),
			Size:          string(product.GetStringBytes("EU")),
			ATCSize:       string(product.GetStringBytes("sizeName")),
			Price:         price,
		}

		products.Products = append(products.Products, buzzObject)

		if buzzObject.Quantity > 0 {
			products.InstockProducts = append(products.InstockProducts, buzzObject)
		}
	}

	if m.Size != "random" {
		for _, product := range products.Products {
			if product.Size == m.Size || product.ATCSize == m.Size {
				logs.Log(m, "Found size: "+m.Size)
				return NOTIFY_TASKS
			}
		}

		logs.Log(m, "Could not find size, retrying...")
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	if loadedSizeQuantity <= 0 {
		logs.Log(m, "Product has loaded stock, however no sizes have pairs loaded yet. ", totalQuantity)
		time.Sleep(m.Delay)
		return GET_PRODUCT_INFO
	}

	b.Products = products

	return NOTIFY_TASKS
}
