package athletesfoot

import "github.com/HellasAIO/HellasAIO/internal/task"

var (
	WAIT_FOR_MONITOR task.TaskState = "wait_for_monitor"
	GET_ORDER_ID     task.TaskState = "get_order_id"
	SUBMIT_ORDER     task.TaskState = "submit_order"
	CHECKOUT_ORDER   task.TaskState = "checkout_order"
)
