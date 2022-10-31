package athletesfoot

import (
	"github.com/HellasAIO/HellasAIO/internal/task"
)

func Initialize() {
	taskType := task.RegisterTaskType("athletesfoot", 0)

	taskType.SetFirstHandlerState(WAIT_FOR_MONITOR)

	taskType.AddHandlers(task.TaskHandlerMap{
		WAIT_FOR_MONITOR: waitForMonitor,
		GET_ORDER_ID:     getOrderId,
		SUBMIT_ORDER:     submitOrder,
		CHECKOUT_ORDER:   checkout,
	})
}
