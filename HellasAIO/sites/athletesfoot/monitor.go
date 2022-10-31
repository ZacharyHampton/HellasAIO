package athletesfoot

import (
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	athletesfootmonitor "github.com/HellasAIO/HellasAIO/monitors/athletesfoot"
)

func waitForMonitor(t *task.Task, i *AthletesFootTaskInternal) task.TaskState {
	if !utils.Contains([]string{"msku"}, t.Mode) {
		logs.Log(t, "Mode is not supported for this site.")
		return task.ErrorTaskState
	}

	t.MonitorData = make(chan interface{})
	logs.Log(t, "Waiting for monitor...")

	for {
		select {
		case data := <-t.MonitorData:
			logs.Log(t, "Got monitor data")
			monitorData := data.(*athletesfootmonitor.AthletesFootMonitorData)
			t.Client = monitorData.Client
			t.CheckoutData = monitorData.CheckoutData

			if t.Size != "random" {
				t.CheckoutData.Size = t.Size
			}

			return GET_ORDER_ID
		}
	}
}
