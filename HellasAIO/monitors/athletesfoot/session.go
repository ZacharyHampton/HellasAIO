package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/task"
	"time"
)

func getSession(m *task.Task, a *AthletesFootInternal) task.TaskState {
	_, err := m.Client.NewRequest().
		SetURL("https://www.theathletesfoot.gr/services/api/Consents/GetCookiesConsentsModel?v=0.3721624200789917&lang=el").
		SetMethod("GET").
		SetDefaultHeadersAF().
		Do()

	if err != nil {
		// handle error and retry
		return GET_SESSION
	}

	// send a request to get session
	return HandleSessionResponse(m, a)
}

func HandleSessionResponse(m *task.Task, a *AthletesFootInternal) task.TaskState {
	if m.Client.LatestResponse.StatusCode() != 200 {
		// retry
		time.Sleep(m.Delay)
		return GET_SESSION
	}

	return LOGIN
}
