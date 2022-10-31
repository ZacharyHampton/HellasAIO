package rich_presence

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/hugolgst/rich-go/client"
)

func SetSite(siteId int) {
	siteName := utils.SiteIDtoSiteStringProper[siteId]
	CurrentActivity.State = "Running " + siteName

	err := client.SetActivity(CurrentActivity)
	if err != nil {
		fmt.Println("Failed to set running site for discord rich presence.")
		return
	}
}
