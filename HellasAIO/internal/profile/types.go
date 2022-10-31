package profile

import "github.com/iancoleman/orderedmap"

type Profile struct {
	Name         string `json:"name"`
	ProfileGroup string `json:"profileGroup"`
	Address      struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		HomePhone   string `json:"homePhone"`
		MobilePhone string `json:"mobilePhone"`
		Address     string `json:"address"`
		ZipCode     string `json:"zipCode"`
		City        string `json:"city"`
		Area        string `json:"area"`
		Prefecture  string `json:"prefecture"`
	} `json:"address"`
}

type ProfileGroup struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Profiles *orderedmap.OrderedMap `json:"profiles"` // ordered map to make sure our profile selection works
}
