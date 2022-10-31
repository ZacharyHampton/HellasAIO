package fuelmonitor

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
)

type FuelInternal struct {
	CartId         string
	ParentSKU      string
	VariantSKU     string
	ProductFound   bool
	Account        *account.Account
	ProfileFound   bool
	Profile        *profile.Profile
	CheckoutData   *logs.CheckoutLogRequest
	TimesFailed    int
	ProxyURL       string
	SessionExisted bool
}

type Variant struct {
	name      string // Nike Air Force 1 '07 LV8 - EU 41
	sku       string // DV3501-400-EU 41
	id        int    // 70351
	sizeId    int    // 5860
	isInStock bool   // data: "IN_STOCK" -> true
}

type FuelMonitorData struct {
	Client *hclient.Client
}

type GraphQLRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

//var _ = reflect.TypeOf(GraphQLRequest{})

type SizesRequestVariables struct {
	MSKU string `json:"msku"`
}

//var _ = reflect.TypeOf(SizesRequestVariables{})

type SizesRequest struct {
	Query     string                `json:"query"`
	Variables SizesRequestVariables `json:"variables"`
}

//var _ = reflect.TypeOf(SizesRequest{})

type CheckoutRequestVariablesAddress struct {
	Firstname         string      `json:"firstname"`
	Lastname          string      `json:"lastname"`
	Company           interface{} `json:"company"`
	Street            []string    `json:"street"`
	City              string      `json:"city"`
	Region            string      `json:"region"`
	Postcode          string      `json:"postcode"`
	CountryCode       string      `json:"country_code"`
	Telephone         string      `json:"telephone"`
	SaveInAddressBook bool        `json:"save_in_address_book"`
}

//var _ = reflect.TypeOf(CheckoutRequestVariablesAddress{})

type CheckoutRequestVariables struct {
	CartId     string                          `json:"cartId"`
	MSKU       string                          `json:"MSKU"`
	VariantSKU string                          `json:"variantSKU"`
	Address    CheckoutRequestVariablesAddress `json:"address"`
}

//var _ = reflect.TypeOf(CheckoutRequestVariables{})

type CheckoutRequest struct {
	Query     string                   `json:"query"`
	Variables CheckoutRequestVariables `json:"variables"`
}

//var _ = reflect.TypeOf(CheckoutRequest{})

func ParentToVariantConversion(parentSKU, size string) string {
	return fmt.Sprintf("%s-EU %s", parentSKU, size)
}

var SizeConversion = map[string]string{
	"18.5":     "5829",
	"19":       "5830",
	"19.5":     "5831",
	"20":       "5832",
	"20.5":     "5833", // assumption
	"21":       "5834",
	"21.5":     "5835", // assumption
	"22":       "5836",
	"22.5":     "5837", // assumption
	"23":       "5838",
	"23.5":     "5839",
	"24":       "5840",
	"25":       "5790",
	"25.5":     "5842",
	"26":       "5791",
	"26.5":     "5843",
	"27":       "5792",
	"27.5":     "5844",
	"28":       "5793",
	"28.5":     "5845",
	"29":       "5794",
	"29.5":     "5846",
	"30":       "5795",
	"31":       "5796",
	"31.5":     "5848",
	"32":       "5797",
	"32.5":     "5849",
	"33":       "5801",
	"33.5":     "5850",
	"34":       "5802",
	"35.5":     "5853",
	"36":       "5727",
	"36-41":    "5717",
	"36.5":     "5854",
	"37":       "5855",
	"37.5":     "5856",
	"38":       "5809",
	"38-42":    "5718",
	"38.5":     "5857",
	"39":       "5858",
	"39.5":     "5859",
	"39-42":    "5728",
	"40":       "5730",
	"40.5":     "5731",
	"41":       "5860",
	"41.5":     "5861",
	"42":       "5732",
	"42.5":     "5862",
	"42-47":    "5921",
	"43":       "5733",
	"43.5":     "5864",
	"43-46":    "5734",
	"44":       "5810",
	"44.5":     "5735",
	"45":       "5823",
	"45.5":     "5736",
	"46":       "5824",
	"46.5":     "5868",
	"47":       "5825",
	"47.5":     "5869",
	"48.5":     "5870",
	"49.5":     "5827",
	"ONE SIZE": "5747",
	"XXS":      "5815",
	"XS":       "5751",
	"S":        "5445",
	"S-M":      "5917",
	"M":        "5446",
	"L":        "5447",
	"L-XL":     "5918",
	"XL":       "5752",
	"XXL":      "5753",
	"XXXL":     "5817",
	"M-L":      "5942",
}
