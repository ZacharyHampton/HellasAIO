package fuelmonitor

import (
	"bytes"
	"github.com/HellasAIO/HellasAIO/internal/cloudflare"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/valyala/fastjson"
	"strings"
	"time"
)

func checkout(m *task.Task, f *FuelInternal) task.TaskState {
	if !f.ProductFound {
		return task.ErrorTaskState
	}

	if !f.ProfileFound {
		f.Profile, _ = profile.GetProfileById(m.ProfileId)
		f.ProfileFound = true
	}

	nameSplit := strings.Split(f.Profile.Address.Name, " ")

	requestBody := CheckoutRequest{
		Query: `mutation checkout($cartId: String!, $variantSKU: String!, $MSKU: String!, $address: CartAddressInput!) {
    addConfigurableProductsToCart(input: {cart_id: $cartId, cart_items: [{parent_sku: $MSKU, data: {quantity: 1, sku: $variantSKU}}]}) {
        cart {
            prices {
                subtotal_including_tax {
                    currency
                    value
                }
            }
            items {
                product {
                    name
                    sku
                    small_image {
                        url
                    }
                    price {
                        regularPrice {
                            amount {
                                currency
                                value
                            }
                        }
                    }
                }
                quantity
            }
        }
    }

    setBillingAddressOnCart(input: {cart_id: $cartId, billing_address: {address: $address, use_for_shipping: true}}) {
        cart {
            billing_address {
                firstname
                lastname
                company
                street
                city
                region {
                    code
                    label
                }
                postcode
                telephone
                country{
                    code
                    label
                }
            }
        }
    }

    setShippingAddressesOnCart(input: { cart_id: $cartId, shipping_addresses: [{address: $address}]}) {
        cart {
            shipping_addresses {
                firstname
                lastname
                company
                street
                city
                region {
                    code
                    label
                }
                postcode
                telephone
                country {
                    code
                    label
                }
            }
        }
    }

    setShippingMethodsOnCart(input: {cart_id: $cartId, shipping_methods: [{carrier_code: "matrixrate", method_code: "matrixrate_8845"}]}) {
        cart {
        shipping_addresses {
            selected_shipping_method {
            carrier_code
            carrier_title
            method_code
            method_title
            amount {
                value
                currency
            }
            }
        }
        }
    }

    setPaymentMethodOnCart(input: {cart_id: $cartId, payment_method: {code: "msp_cashondelivery"}}) {
        cart {
            selected_payment_method {
                code
                title
            }
        }
    }

    placeOrder(input: {cart_id: $cartId}) {
        order {
            order_id
        }
    }
}`,
		Variables: CheckoutRequestVariables{CartId: f.CartId, MSKU: f.ParentSKU, VariantSKU: f.VariantSKU, Address: CheckoutRequestVariablesAddress{
			Firstname:         nameSplit[0],
			Lastname:          nameSplit[1],
			Company:           nil,
			Street:            []string{f.Profile.Address.Address},
			City:              f.Profile.Address.City,
			Region:            f.Profile.Address.Prefecture,
			Postcode:          f.Profile.Address.ZipCode,
			CountryCode:       "GR",
			Telephone:         f.Profile.Address.MobilePhone,
			SaveInAddressBook: false,
		}},
	}

	_, err := m.Client.NewRequest().
		SetURL("https://www.fuel.com.gr/el/graphql").
		SetMethod("POST").
		SetHeader("User-Agent", utils.UserAgent).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "*/*").
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		return CHECKOUT
	}

	return handleCheckoutResponse(m, f)
}

func handleCheckoutResponse(m *task.Task, f *FuelInternal) task.TaskState {
	if cloudflare.DetectCloudflare(m) {
		logs.Log(m, "Cloudflare detected.")
		for {
			logs.Log(m, "Solving cloudflare...")
			success := cloudflare.GetClearanceCookie(m, "https://www.fuel.com.gr", f.ProxyURL)
			if success {
				logs.Log(m, "Cloudflare solved.")
				return GET_CART_ID
			} else {
				logs.Log(m, "Failed to solve cloudflare. Retrying...")
				time.Sleep(m.Delay)
			}
		}
	}

	//m.Log("Checkout response: ", m.Client.LatestResponse.BodyAsString())
	if bytes.Contains(m.Client.LatestResponse.Body(), []byte("errors")) {
		jData, _ := fastjson.ParseBytes(m.Client.LatestResponse.Body())
		for _, jValues := range jData.GetArray("errors") {
			if string(jValues.GetStringBytes("path", "0")) == "addConfigurableProductsToCart" {
				if bytes.Contains(jValues.GetStringBytes("message"), []byte("Could not find specified product.")) {
					logs.Log(m, "Product variant not found")
					return CHECKOUT
				}

				if bytes.Contains(jValues.GetStringBytes("message"), []byte("This product is out of stock.")) {
					logs.Log(m, "Product OOS")
					return CHECKOUT
				}

				if bytes.Contains(jValues.GetStringBytes("message"), []byte("Could not add the product with SKU")) {
					logs.Log(m, "Product not found.")
					return CHECKOUT
				}

			}

			logs.Log(m, "Checkout error: ", string(jValues.GetStringBytes("message")))
		}

		logs.Log(m, "Checkout failed.")
		title.AddFailure()
		f.TimesFailed += 1
		if f.TimesFailed < 3 {
			return GET_CART_ID
		}

		m.CheckoutData.Status = "denied"
		return CHECKOUT

	} else {
		logs.Log(m, "Checkout succeeded!")
		m.CheckoutData.Status = "success"
		title.AddCheckout()
		title.AddCart()
	}

	m.CheckoutData.Price = fastjson.GetFloat64(m.Client.LatestResponse.Body(), "data", "addConfigurableProductsToCart", "cart", "prices", "subtotal_including_tax", "value")
	m.CheckoutData.ProductName = fastjson.GetString(m.Client.LatestResponse.Body(), "data", "addConfigurableProductsToCart", "cart", "items", "0", "product", "name")
	m.CheckoutData.ImageUrl = fastjson.GetString(m.Client.LatestResponse.Body(), "data", "addConfigurableProductsToCart", "cart", "items", "0", "product", "small_image", "url")

	return task.DoneTaskState
}
