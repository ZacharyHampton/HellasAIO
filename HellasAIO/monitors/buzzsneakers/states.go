package buzzsneakersmonitor

import "github.com/HellasAIO/HellasAIO/internal/task"

var (
	INITIALIZE       task.TaskState = "initialize"
	GET_ITEM         task.TaskState = "get_item"
	GET_PRODUCT_INFO task.TaskState = "get_product_info"
	NOTIFY_TASKS     task.TaskState = "notify_tasks"
)
