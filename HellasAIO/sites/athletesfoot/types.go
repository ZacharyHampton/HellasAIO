package athletesfoot

import "github.com/HellasAIO/HellasAIO/internal/profile"

var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

type AthletesFootTaskInternal struct {
	OrderId          int
	ProfileRetrieved bool
	Profile          *profile.Profile
}

type SubmitOrderRequest struct {
	Data      string `json:"data"`
	Arguments string `json:"arguments"`
}
