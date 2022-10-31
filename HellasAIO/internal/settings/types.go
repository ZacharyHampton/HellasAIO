package settings

type Settings struct {
	AuthKey            string  `json:"key"`
	DiscordWebhook     string  `json:"discord_webhook"`
	AllowPublicWebhook bool    `json:"allowPublicWebhook"`
	Captcha            Captcha `json:"captcha"`
}

type Captcha struct {
	TwoCaptcha CaptchaDetails `json:"2captcha"`
	CapMonster CaptchaDetails `json:"capmonster"`
}

type CaptchaDetails struct {
	Key string `json:"key"`
}
