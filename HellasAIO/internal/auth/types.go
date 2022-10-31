package auth

type AuthRequest struct {
	Key  string `json:"licenseKey"`
	HWID string `json:"HWID"`
}
