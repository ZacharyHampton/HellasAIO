package buzzsneakersmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/task"
)

type BuzzSneakersInternal struct {
	Account   *account.Account
	ProductId string
	Products  BuzzProduct
}

type BuzzMonitorData struct {
	CheckoutData task.CheckoutLogRequest
	Products     BuzzProduct
	ProductId    string
}

type BuzzCombinationProduct struct {
	Quantity      float64
	CombinationId string
	Size          string
	ATCSize       string
	Price         float64
}

type BuzzProduct struct {
	InstockProducts []BuzzCombinationProduct
	Products        []BuzzCombinationProduct
}
