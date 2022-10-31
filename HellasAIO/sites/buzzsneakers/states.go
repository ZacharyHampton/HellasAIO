package buzzsneakers

import "github.com/HellasAIO/HellasAIO/internal/task"

var (
	INITIALIZE       task.TaskState = "initialize"
	LOGIN            task.TaskState = "login"
	WAIT_FOR_MONITOR task.TaskState = "wait_for_monitor"
	ADD_TO_CART      task.TaskState = "add_to_cart"
	CHECKOUT_ORDER   task.TaskState = "checkout_order"
)
