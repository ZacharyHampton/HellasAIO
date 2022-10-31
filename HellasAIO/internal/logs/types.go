package logs

import (
	"time"
)

type CheckoutLogRequest struct {
	TaskStart   time.Time `json:"-"`
	TaskEnd     time.Time `json:"-"`
	Price       float64   `json:"price"`
	ProductName string    `json:"product_name"`
	ProductMSKU string    `json:"product_msku"`
	Mode        string    `json:"mode"`
	CheckoutMs  int       `json:"checkout_ms"`
	Size        string    `json:"size"`
	Status      string    `json:"status"`
	Website     string    `json:"website"`
	ImageUrl    string    `json:"image_url"`
	AllowPublic bool      `json:"allow_public"`
}

type LogtailData struct {
	AuthKey     string `json:"auth_key"`
	SiteName    string `json:"site_name"`
	TaskType    string `json:"task_type"`
	TaskMode    string `json:"task_mode"`
	TaskProduct string `json:"task_product"`
	Version     string `json:"version"`
	Message     string `json:"message"`
	Count       int    `json:"count"`
}
