package athletesfoot

import (
	"bytes"
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"strconv"
	"strings"
	"time"
)

func submitOrder(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	if !i.ProfileRetrieved {
		profileObject, err := profile.GetProfileById(t.ProfileId)

		if err != nil {
			// handle error and retry
			return SUBMIT_ORDER
		}

		i.Profile = profileObject
		i.ProfileRetrieved = true
	}

	nSplit := strings.Split(i.Profile.Address.Name, " ")
	FirstName, LastName := nSplit[0], nSplit[1]

	sOrderID := strconv.Itoa(i.OrderId)
	requestBody := SubmitOrderRequest{
		Data: fmt.Sprintf("{\"ControllerName\":\"CheckoutOrderController\",\"OrderId\":\"%s\",\"ProductsSource\":0,\"PhotoGalleryImageTypeCode\":null}", sOrderID),
		Arguments: fmt.Sprintf(
			"[ { \"FirstName\": \"\", \"LastName\": \"\", \"UserName\": \"\", \"Email\": \"%s\", \"Gender\": \"\", \"Phone\": \"\", \"Mobile\": \"\", \"Newsletter\": false, \"LoyaltyMemberCardNumber\": \"\", \"SocialSecurityNumber\": \"\", \"CreateAccount\": false, \"UseEmailAsLoginName\": false, \"HashedEmail\": \"%s\" }, { \"Id\": 1, \"Code\": \"onSite\", \"Name\": \"a\", \"Description\": \"b\", \"ShippingMethods\": [ 1 ] }, { \"Id\": 0, \"Name\": \"\", \"Description\": \"\", \"Address\": \"\", \"Area\": \"\", \"City\": \"\", \"PostCode\": \"\", \"Email\": \"\", \"Phone\": \"\", \"Mobile\": \"\", \"Fax\": \"\", \"Contact\": \"\", \"OpeningHours\": \"\", \"ClosingHours\": \"\", \"deliveryTime\": \"\", \"Latitude\": \"\", \"Longitude\": \"\", \"Country\": \"\", \"CountryCode\": \"\" }, { \"Id\": 36750, \"Title\": \"\", \"FirstName\": \"%s\", \"LastName\": \"%s\", \"Email\": \"%s\", \"Address\": \"%s\", \"Address2\": \"\", \"Area\": \"%s\", \"City\": \"%s\", \"Prefecture\": \"%s\", \"PostCode\": \"%s\", \"Country\": \"GR\", \"Phone\": \"%s\", \"Mobile\": \"%s\", \"Fax\": \"\", \"Bell\": \"\", \"Comments\": \"\", \"Floor\": \"\", \"Apartment\": \"\", \"Block\": \"\", \"Entrance\": \"\", \"Streetnumber\": \"\", \"GeoZones\": [ 14, 20 ] }, \"\", { \"Id\": 36750, \"Title\": \"\", \"FirstName\": \"%s\", \"LastName\": \"%s\", \"Email\": \"%s\", \"Address\": \"%s\", \"Address2\": \"\", \"Area\": \"%s\", \"City\": \"%s\", \"Prefecture\": \"%s\", \"PostCode\": \"%s\", \"Country\": \"GR\", \"Phone\": \"%s\", \"Mobile\": \"%s\", \"Fax\": \"\", \"Bell\": \"\", \"Comments\": \"\", \"Floor\": \"\", \"Apartment\": \"\", \"Block\": \"\", \"Entrance\": \"\", \"Streetnumber\": \"\", \"GeoZones\": [ 14, 20 ] }, { \"MethodId\": 1, \"MethodName\": \"c\", \"MethodDescription\": \"\", \"CompanyId\": 9, \"Duration\": \"4-10\", \"AdditionalDuration\": \"\", \"DurationText\": \"d\", \"AdditionalDurationText\": \"\", \"CompanyName\": \"Speedex\", \"CompanyPhones\": \"\", \"CompanyAddress\": \"\", \"UserCompany\": { \"Name\": null, \"Address\": null, \"Town\": null, \"Area\": null, \"Phone\": null, \"PostalCode\": null }, \"PartnerDurationText\": \"\", \"FromStoreTimeText\": \"\", \"PaymentMethods\": [ 3, 1, 5 ], \"AvailableShippingCompanies\": [ 17, 18, 9 ] }, { \"Name\": \"e\", \"Code\": \"cashOnDelivery\", \"MethodId\": 3, \"ProcessorId\": 5, \"InstallmentsCount\": 0 }, null, { }, { \"Id\": 1, \"IsVisible\": true, \"Code\": \"Normal\", \"Description\": \"f\", \"DescriptionEl\": \"g\", \"DescriptionEn\": \"Normal\", \"IsPaymentOrderType\": false }, \"\", \"\", { }, \"\", { \"ExplicitTimeslotId\": 0, \"DeliveryPointId\": 0, \"Name\": \"\", \"DateTimeFrom\": \"0001-01-01T00:00:00\", \"DateTimeTo\": \"0001-01-01T00:00:00\", \"EstimatedDeliveryDateTimeFrom\": \"0001-01-01T00:00:00\", \"EstimatedDeliveryDateTimeTo\": \"0001-01-01T00:00:00\", \"MaxUsages\": 0, \"CurrentCount\": 0, \"IsActive\": false, \"AdditionalCharge\": 0 }, null]",
			i.Profile.Address.Email, utils.GetMD5Hash(i.Profile.Address.Email), FirstName, LastName, i.Profile.Address.Email, i.Profile.Address.Address, i.Profile.Address.Area, i.Profile.Address.City, i.Profile.Address.Prefecture, i.Profile.Address.ZipCode, i.Profile.Address.HomePhone, i.Profile.Address.MobilePhone, FirstName, LastName, i.Profile.Address.Email, i.Profile.Address.Address, i.Profile.Address.Area, i.Profile.Address.City, i.Profile.Address.Prefecture, i.Profile.Address.ZipCode, i.Profile.Address.HomePhone, i.Profile.Address.MobilePhone,
		),
	}

	_, err := t.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/services/checkoutservice.svc/SubmitOrder?lang=el").
		SetMethod("POST").
		SetDefaultHeadersAF().
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		logs.Log(t, "Error submitting order.")
		time.Sleep(t.Delay)
		return SUBMIT_ORDER
	}

	return handleSubmitOrderResponse(t, i)
}

func handleSubmitOrderResponse(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	if !bytes.Contains(t.Client.LatestResponse.Body(), []byte("OrderSubmit_OperationSuccessful")) {
		logs.Log(t, "Failed to submit order.")
		time.Sleep(t.Delay)
		return SUBMIT_ORDER
	}

	return CHECKOUT_ORDER
}
