package rich_presence

import (
	"fmt"
	"github.com/hugolgst/rich-go/client"
	"time"
)

var (
	startingTime    = time.Now()
	CurrentActivity = client.Activity{
		State:      "Browsing Menu",
		Details:    "Destroying Greek Websites",
		LargeImage: "hellasaiologo",
		LargeText:  "HellasAIO",
		SmallImage: "hellasaiologo",
		SmallText:  "HellasAIO",
		Party: &client.Party{
			ID:         "-1",
			Players:    1,
			MaxPlayers: 1,
		},
		Timestamps: &client.Timestamps{
			Start: &startingTime,
		},
		Buttons: []*client.Button{
			{
				Label: "Twitter",
				Url:   "https://twitter.com/hellasaio",
			},
			{
				Label: "Website",
				Url:   "https://hellasaio.com/",
			},
		},
	}
)

func Initialize() {
	err := client.Login("discordid")
	if err != nil {
		fmt.Println("Failed to start discord rich presence.")
		return
	}

	err = client.SetActivity(CurrentActivity)
	if err != nil {
		fmt.Println("Failed to set HellasAIO discord rich presence.")
		return
	}
}
