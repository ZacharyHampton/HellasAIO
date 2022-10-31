package task

import (
	"context"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
	"reflect"
	"time"
)

type Task struct {
	TaskType     int                // taskType ? 0:1 | 0 = Monitor | 1 = Checkout
	SiteId       int                // siteId
	ID           string             `json:"id"`          // task uuid
	Type         string             `json:"type"`        // registered task type aka site name
	Mode         string             `json:"mode"`        // task mode
	Product      string             `json:"product"`     // product info
	Size         string             `json:"size"`        // size info
	ProxyListID  string             `json:"proxyListID"` // proxy list id
	ProfileId    string             `json:"profileId"`   // profile id
	AccountId    string             `json:"accountId"`   // account id
	Delay        time.Duration      `json:"delay"`       // delay (in ms)
	Internal     interface{}        `json:"-"`           // internal data, gotten from second func argument
	Active       bool               `json:"-"`           // active status
	MonitorData  chan interface{}   `json:"-"`           // monitor data, only used in checkout tasks
	Context      context.Context    `json:"-"`
	Cancel       context.CancelFunc `json:"-"` // cancel function
	Client       *hclient.Client    `json:"-"` // http client
	CheckoutData CheckoutLogRequest `json:"-"` // checkout data
}

type CheckoutLogRequest struct {
	TaskStart   time.Time `json:"-"`            // auto defined
	TaskEnd     time.Time `json:"-"`            // auto defined
	CheckoutMs  int       `json:"checkout_ms"`  // auto defined
	Price       float64   `json:"price"`        // needs to be defined
	ProductName string    `json:"product_name"` // needs to be defined
	ProductMSKU string    `json:"product_msku"` // needs to be defined
	Mode        string    `json:"mode"`         // needs to be defined
	Size        string    `json:"size"`         // needs to be defined
	Status      string    `json:"status"`       // needs to be defined
	Website     string    `json:"website"`      // siteName, needs to be defined
	ImageUrl    string    `json:"image_url"`    // needs to be defined
}

// taskType ? 0:1 | 0 = Monitor | 1 = Checkout

type TaskGroup struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Tasks map[string]bool `json:"task"`
}

type TaskType struct {
	firstHandlerState TaskState
	internalType      reflect.Type
	handlers          TaskReflectMap
}

type TaskState string
type TaskHandlerMap map[TaskState]interface{}
type TaskReflectMap map[string]reflect.Value
