package sessions

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"os"
	"strings"
)

var (
	Paths = []string{
		"../.sessions/athletesfoot",
		"../.sessions/fuel",
		"../.sessions/slamdunk",
		"../.sessions/buzzsneakers",
		"../.sessions/europesports",
	}
)

func Initialize() {
	for _, path := range Paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
				return
			}
		}
	}
}

func DoesSessionExist(account *account.Account) bool {
	path := fmt.Sprintf("../.sessions/%s/%s.sessions", strings.Replace(utils.SiteIDtoSiteString[account.SiteId], "@", "", -1), account.Email)
	_, err := os.Stat(path)
	return err == nil
}
