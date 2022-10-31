package loading

import (
	"github.com/HellasAIO/HellasAIO/internal/account"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/settings"
)

type Accounts struct {
	Accounts map[int][]account.Account
}

type Proxies struct {
	Proxies []proxy.ProxyGroup
}

type Tasks struct {
	Tasks map[int][]string
}

type Profiles struct {
	Profiles []profile.Profile
}

type Settings struct {
	Settings settings.Settings `json:"settings"`
}

type QuicktaskGroup struct {
	ProfileName  string
	AccountEmail string
}

type Config struct {
	Accounts        Accounts
	Proxies         Proxies
	Tasks           Tasks
	Profiles        Profiles
	QuicktaskGroups map[int][]QuicktaskGroup
	Settings        Settings `json:"settings"`
}

var TaskModeAndSiteIDToRegisteredSiteName = map[string]string{
	"monitor,0":  "athletesfootmonitor",
	"checkout,0": "athletesfoot",
	"monitor,1":  "fuelmonitor",
	"checkout,1": "fuelmonitor",
	"monitor,2":  "slamdunkmonitor",
	"checkout,2": "slamdunk",
	"monitor,3":  "buzzsneakersmonitor",
	"checkout,3": "buzzsneakers",
	"monitor,4":  "europesportsmonitor",
	"checkout,4": "europesports",
}
