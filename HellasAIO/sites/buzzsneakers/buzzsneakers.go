package buzzsneakers

import "github.com/HellasAIO/HellasAIO/internal/task"

func Initialize() {
	taskType := task.RegisterTaskType("buzzsneakers", 3)

	taskType.SetFirstHandlerState(INITIALIZE)

	taskType.AddHandlers(task.TaskHandlerMap{
		INITIALIZE:       initialize,
		LOGIN:            login,
		WAIT_FOR_MONITOR: waitForMonitor,
		ADD_TO_CART:      addToCart,
		CHECKOUT_ORDER:   checkout,
	})
}
