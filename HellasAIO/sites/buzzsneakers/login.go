package buzzsneakers

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/valyala/fastjson"
	"net/url"
	"time"
)

func login(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	requestBody := url.Values{}
	requestBody.Add("login_email", b.Account.Email)
	requestBody.Add("login_password", b.Account.Password)
	requestBody.Add("back_url", "https://www.buzzsneakers.gr/oloklirosi-parangelias")
	requestBody.Add("ajax", "yes")
	requestBody.Add("task", "login")

	_, err := c.Client.NewRequest().
		SetURL("https://www.buzzsneakers.gr/eisodos").
		SetMethod("POST").
		SetDefaultHeadersBuzz().
		SetFormBody(requestBody).
		Do()

	if err != nil {
		logs.Log(c, "Error logging in.")
		time.Sleep(c.Delay)
		return LOGIN
	}

	return handleLoginResponse(c, b)
}

func handleLoginResponse(c *task.Task, b *BuzzCheckoutInternal) task.TaskState {
	if fastjson.GetBool(c.Client.LatestResponse.Body(), "flag") {
		logs.Log(c, "Successfully logged in.")
		c.Client.SaveCookies()

		if c.Mode == "login" {
			return task.DoneTaskState
		}

		return WAIT_FOR_MONITOR
	} else {
		logs.Log(c, "Failed to login.")
		return task.ErrorTaskState
	}
}
