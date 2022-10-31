package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
)

type LoginRequest struct {
	Email         string `json:"loginName"`
	Password      string `json:"password"`
	RememberMe    bool   `json:"rememberMe"`
	ReCaptchaCode string `json:"reCaptchaCode"`
	AccessToken   string `json:"accessToken"`
}

type GetCartRequest struct {
	CurrentItemId int    `json:"currentItemId"`
	TemplateCode  string `json:"templateCode"`
}

type RemoveFromCartRequest struct {
	Id int `json:"id"`
}

type AthletesFootInternal struct {
	Account   *account.Account
	foundItem bool
	sku       int
	subSku    int
	pathId    int
}

type GetItemRequest struct {
	EsFields    string `json:"es_fields"`
	EsGeoWindow string `json:"es_geo_window"`
	PageFrom    int    `json:"pagefrom"`
	PageSize    int    `json:"pagesize"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Type        string `json:"type"`
}

type AddToCartRequest struct {
	Item struct {
		SKU                               int      `json:"id"`
		SubSKU                            int      `json:"skuId"`
		PathId                            int      `json:"path"`
		Quantity                          string   `json:"qty"`
		ExtraAttributes                   string   `json:"extraAttributes"`
		BundleItems                       []string `json:"bundleItems"`
		RecipeItems                       []string `json:"recipeItems"`
		ComboItems                        []string `json:"comboItems"`
		IsDefaultRecipeQuantitiesIncluded bool     `json:"isDefaultRecipeQuantitiesIncluded"`
		EnhancedInfo                      struct {
			MinQuantity int `json:"minQuantity"`
			MaxQuantity int `json:"maxQuantity"`
			Quanta      int `json:"quanta"`
		} `json:"enhancedInfo"`
	} `json:"item"`
}

type AthletesFootMonitorData struct {
	Client       *hclient.Client
	CheckoutData task.CheckoutLogRequest
}
