package account

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	SiteId   int    `json:"siteId"`
}
