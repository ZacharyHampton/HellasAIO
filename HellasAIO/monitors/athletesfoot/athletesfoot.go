package athletesfootmonitor

import (
	"github.com/HellasAIO/HellasAIO/internal/task"
)

func Initialize() {
	monitorType := task.RegisterTaskType("athletesfootmonitor", 0)

	monitorType.SetFirstHandlerState(INITIALIZE)

	monitorType.AddHandlers(task.TaskHandlerMap{
		INITIALIZE:  initialize,
		GET_SESSION: getSession,
		LOGIN:       login,
		CLEAR_CART:  clearCart,
		GET_ITEM:    getItem,
		ADD_TO_CART: addToCart,
	})
}
