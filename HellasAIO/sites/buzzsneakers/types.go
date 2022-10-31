package buzzsneakers

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	buzzsneakersmonitor "github.com/HellasAIO/HellasAIO/monitors/buzzsneakers"
)

type BuzzCheckoutInternal struct {
	ProfileFound bool
	Profile      *profile.Profile
	Account      *account.Account
	Products     buzzsneakersmonitor.BuzzProduct
	Combination  *buzzsneakersmonitor.BuzzCombinationProduct
	ProductID    string
}
