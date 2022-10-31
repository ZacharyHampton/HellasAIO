package fuelmonitor

import "github.com/HellasAIO/HellasAIO/internal/task"

func Initialize() {
	monitorType := task.RegisterTaskType("fuelmonitor", 1)

	monitorType.SetFirstHandlerState(INITIALIZE)

	monitorType.AddHandlers(task.TaskHandlerMap{
		INITIALIZE:       initialize,
		SOLVE_CLOUDFLARE: solveCloudflare,
		LOGIN:            login,
		GET_SIZE:         findSizes,
		GET_CART_ID:      getCartId,
		CHECKOUT:         checkout,
	})
}
