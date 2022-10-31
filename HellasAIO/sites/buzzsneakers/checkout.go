package buzzsneakers

import (
	"bytes"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"net/url"
	"strings"
	"time"
)

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func checkout(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	/*
		quantity_1=1&cart_comment=&ticket_nb_action=&ua=1&cart_ticket_type=nb_action&ticket=&cart_onepage_type_person=1&cart_onepage_firstname=John&cart_onepage_lastname=Ioannou&cart_onepage_email=zhampton77%40gmail.com&cart_onepage_phone=6969696969&cart_onepage_region_id=70&cart_onepage_city=%CE%92%CE%9F%CE%9B%CE%9F%CE%A3+78758&cart_onepage_city_id=10621&cart_onepage_postcode=78758&cart_onepage_street=agias+sofias&cart_onepage_street_id=0&cart_onepage_street_no=347&p_firstname=John&p_lastname=Ioannou&p_street=agias+sofias&p_street_no=347&p_postcode=78758&p_city=%CE%92%CE%9F%CE%9B%CE%9F%CE%A3&p_phone=6969696969&orderAddress=yes&carierId=10&cart_onepage_deliveryTime_10=-1&typepayment=post&cart_onepage_terms_of_use=1&submit_order_one_page=1
	*/

	if !b.ProfileFound {
		b.Profile, _ = profile.GetProfileById(c.ProfileId)
		b.ProfileFound = true
	}

	nameSplit := strings.Split(b.Profile.Address.Name, " ")
	streetSplit := strings.Split(b.Profile.Address.Address, " ")
	streetNumbers := streetSplit[len(streetSplit)-1]
	streetSplit = remove(streetSplit, len(streetSplit)-1)
	street := strings.Join(streetSplit, " ")

	requestBody := url.Values{}
	requestBody.Set("quantity_1", "1")
	requestBody.Set("cart_comment", "")
	requestBody.Set("ticket_nb_action", "")
	requestBody.Set("ua", "1")
	requestBody.Set("cart_ticket_type", "nb_action")
	requestBody.Set("ticket", "")
	requestBody.Set("submit_order_one_page", "1")
	requestBody.Set("cart_onepage_terms_of_use", "1")
	requestBody.Set("typepayment", "post")
	requestBody.Set("cart_onepage_deliveryTime_10", "-1")
	requestBody.Set("carierId", "10")
	requestBody.Set("orderAddress", "yes")

	requestBody.Set("cart_onepage_type_person", "1")
	requestBody.Set("cart_onepage_firstname", nameSplit[0])
	requestBody.Set("cart_onepage_lastname", nameSplit[1])
	requestBody.Set("cart_onepage_email", c.AccountId)
	requestBody.Set("cart_onepage_phone", b.Profile.Address.HomePhone)
	requestBody.Set("cart_onepage_city", b.Profile.Address.City)
	requestBody.Set("cart_onepage_postcode", b.Profile.Address.ZipCode)
	requestBody.Set("cart_onepage_street", street)
	requestBody.Set("cart_onepage_street_id", "-1")
	requestBody.Set("cart_onepage_street_no", streetNumbers)

	requestBody.Set("p_firstname", nameSplit[0])
	requestBody.Set("p_lastname", nameSplit[1])
	requestBody.Set("p_street", street)
	requestBody.Set("p_street_no", streetNumbers)
	requestBody.Set("p_postcode", b.Profile.Address.ZipCode)
	requestBody.Set("p_city", b.Profile.Address.City)
	requestBody.Set("p_phone", b.Profile.Address.HomePhone)

	logs.Log(c, "Checking out...")

	_, err := c.Client.NewRequest().
		SetURL("https://www.buzzsneakers.gr/oloklirosi-parangelias").
		SetMethod("POST").
		SetDefaultHeadersBuzz().
		SetFormBody(requestBody).
		Do()

	if err != nil {
		logs.Log(c, "Failed to make checkout request.")
		time.Sleep(c.Delay)
		return CHECKOUT_ORDER
	}

	return handleCheckoutResponse(c, b)
}

func handleCheckoutResponse(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	logs.Log(c, "Got checkout response.")

	if bytes.Contains(c.Client.LatestResponse.Body(), []byte("Undefined index: g-recaptcha-response")) {
		logs.Log(c, "Session expired (captcha flag). Please refresh session.")
		title.AddFailure()
		return task.ErrorTaskState
	}

	if c.Client.LatestResponse.StatusCode() == 200 && bytes.Contains(c.Client.LatestResponse.Body(), []byte("cart-description confirm-info")) {
		logs.Log(c, "Checkout successful.")
		title.AddCheckout()
		c.CheckoutData.Status = "success"
		return task.DoneTaskState
	} else {
		logs.Log(c, "Failed to checkout.")
		title.AddFailure()
		time.Sleep(c.Delay)
		return CHECKOUT_ORDER
	}
}
