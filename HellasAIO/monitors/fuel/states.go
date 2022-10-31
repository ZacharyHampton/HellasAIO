package fuelmonitor

import "github.com/HellasAIO/HellasAIO/internal/task"

var (
	INITIALIZE       task.TaskState = "initialize"
	SOLVE_CLOUDFLARE task.TaskState = "solve_cloudflare"
	LOGIN            task.TaskState = "login"
	GET_CART_ID      task.TaskState = "get_cart_id"
	GET_SIZE         task.TaskState = "get_size"
	CHECKOUT         task.TaskState = "checkout"
)
