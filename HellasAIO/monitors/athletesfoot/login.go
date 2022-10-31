package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"time"
)

func login(m *task.Task, a *AthletesFootInternal) task.TaskState {
	requestBody := LoginRequest{
		Email:         a.Account.Email,
		Password:      a.Account.Password,
		RememberMe:    true,
		ReCaptchaCode: "",
		AccessToken:   "",
	}

	_, err := m.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/services/AuthService.svc/Login?lang=en").
		SetMethod("POST").
		SetDefaultHeadersAF().
		SetJSONBody(requestBody).
		Do()

	if err != nil {
		// handle error and retry
		return LOGIN
	}
	// send a request to
	return HandleLoginResponse(m, a)
}

func HandleLoginResponse(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if !(fastjson.GetString(m.Client.LatestResponse.Body(), "d", "code") == "Login_OperationSuccesful") { // can be faster if you solely deal with bytes, however this is easier to write
		// handle error
		logs.Log(m, "Error while logging in")
		time.Sleep(m.Delay)
		return LOGIN
	}

	// if logged in, run clear cart else return to login
	m.Client.SaveCookies()
	if m.Mode == "login" {
		return task.DoneTaskState
	}

	return CLEAR_CART
}
