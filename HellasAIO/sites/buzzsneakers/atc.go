package buzzsneakers

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"github.com/valyala/fastjson"
	"math/rand"
	"net/url"
	"time"
)

var COUNT = 0

func addToCart(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	products := b.Products.Products

	if c.Size != "random" {
		for _, combination := range products {
			if combination.Size == c.Size {
				b.Combination = &combination
				break
			}
		}
	} else {
		if len(b.Products.InstockProducts) > 0 {
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(b.Products.InstockProducts))
			b.Combination = &b.Products.InstockProducts[randomIndex]
		} else {
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(b.Products.Products))
			b.Combination = &b.Products.Products[randomIndex]
		}
	}

	logs.Log(c, fmt.Sprintf("Adding to cart (size %s)...", b.Combination.ATCSize))

	c.CheckoutData.Size = b.Combination.ATCSize
	c.CheckoutData.Price = b.Combination.Price

	requestBody := url.Values{}
	requestBody.Add("ajax", "yes")
	requestBody.Add("task", "cartInsert")
	requestBody.Add("id", b.ProductID)
	requestBody.Add("combId", b.Combination.CombinationId)
	requestBody.Add("amount", "1")
	requestBody.Add("size", b.Combination.ATCSize)

	_, err := c.Client.NewRequest().
		SetURL("https://www.buzzsneakers.gr/oloklirosi-parangelias").
		SetMethod("POST").
		SetDefaultHeadersBuzz().
		SetHeader("Referrer", "https://www.buzzsneakers.gr/athlitika-papoutsia/"+b.ProductID).
		SetFormBody(requestBody).
		Do()

	if err != nil {
		logs.Log(c, "Error adding to cart.")
		time.Sleep(c.Delay)
		return ADD_TO_CART
	}

	return handleAddToCartRequest(c, b)
}

func handleAddToCartRequest(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	if COUNT >= 20 {
		logs.Log(c, "Failed to ATC 20 times, waiting for new product data.")
		COUNT = 0
		return WAIT_FOR_MONITOR
	}

	if c.Client.LatestResponse.StatusCode() != 200 {
		logs.Log(c, "Error adding to cart.")
		time.Sleep(c.Delay)
		COUNT += 1
		return ADD_TO_CART
	}

	if fastjson.GetBool(c.Client.LatestResponse.Body(), "flag") == false {
		logs.Log(c, "Error adding to cart (OOS or unknown error)")
		time.Sleep(c.Delay)
		COUNT += 1
		return ADD_TO_CART
	}

	logs.Log(c, "Added to cart.")
	title.AddCart()

	return CHECKOUT_ORDER
}
