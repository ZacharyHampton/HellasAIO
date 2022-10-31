package buzzsneakersmonitor

import "github.com/HellasAIO/HellasAIO/internal/task"

func Initialize() {
	monitorType := task.RegisterTaskType("buzzsneakersmonitor", 3)

	monitorType.SetFirstHandlerState(INITIALIZE)

	monitorType.AddHandlers(task.TaskHandlerMap{
		INITIALIZE:       initialize,
		GET_ITEM:         getItem,
		GET_PRODUCT_INFO: getProductInfo,
		NOTIFY_TASKS:     notifyTasks,
	})
}
