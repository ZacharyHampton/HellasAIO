package taskmngr

import (
	"context"
	"github.com/HellasAIO/HellasAIO/internal/loading"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/task"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/getsentry/sentry-go"
	"log"
	"reflect"
	"time"
)

func handleTaskState(taskState task.TaskState, taskType *task.TaskType, t *task.Task) task.TaskState {
	nextTaskHandler, err := taskType.GetHandler(taskState)

	if err != nil {
		log.Println("Task handler error: ", err)
		return task.ErrorTaskState
	}

	// func (t *task.Monitor, internal *SiteInternal) task.TaskState
	return task.TaskState(nextTaskHandler.Call([]reflect.Value{reflect.ValueOf(t), reflect.ValueOf(t.Internal)})[0].String())
}

// RunTask starts a task
func RunTask(t *task.Task) {
	t.Context, t.Cancel = context.WithCancel(context.Background())
	t.Active = true

	defer func() {
		if r := recover(); r != nil {
			log.Println("Task error:", r)

			sentry.RecoverWithContext(t.Context)
			sentry.Flush(time.Second * 5)
		}
	}()

	if !task.DoesTaskTypeExist(t.Type) {
		return
	}

	taskType, err := task.GetTaskType(t.Type)

	if err != nil {
		log.Println("Task type error: ", err)
		t.Active = false
		return
	}

	hasHandlers := taskType.HasHandlers()

	if !hasHandlers {
		t.Active = false
		return
	}

	nextState := taskType.GetFirstHandlerState()
	logs.Log(t, "Starting task...")
	t.CheckoutData.TaskStart = time.Now()

	if len(nextState) == 0 {
		t.Active = false
		return
	}

	t.Internal = reflect.New(taskType.GetInternalType().Elem()).Interface()

	// loop the task states
	for {
		nextState = handleTaskState(nextState, taskType, t)
		if utils.Debug {
			logs.Log(t, nextState)
		}

		if nextState == task.DoneTaskState || t.Context.Err() != nil {
			t.CheckoutData.TaskEnd = time.Now()
			t.CheckoutData.CheckoutMs = int(t.CheckoutData.TaskEnd.Sub(t.CheckoutData.TaskStart).Milliseconds())
			logs.LogCheckout(&logs.CheckoutLogRequest{
				TaskStart:   t.CheckoutData.TaskStart,
				TaskEnd:     t.CheckoutData.TaskEnd,
				Price:       t.CheckoutData.Price,
				ProductName: t.CheckoutData.ProductName,
				ProductMSKU: t.CheckoutData.ProductMSKU,
				Mode:        t.CheckoutData.Mode,
				CheckoutMs:  t.CheckoutData.CheckoutMs,
				Size:        t.CheckoutData.Size,
				Status:      t.CheckoutData.Status,
				Website:     t.CheckoutData.Website,
				ImageUrl:    t.CheckoutData.ImageUrl,
			}, loading.Data.Settings.Settings.DiscordWebhook)
			// you can report that the task stopped here
			t.Active = false
			break
		} else if nextState == task.ErrorTaskState {
			// report errors
			t.Active = false
			break
		}

		time.Sleep(1 * time.Millisecond)
	}
}

// StopTask stops a task
func StopTask(t *task.Task) {
	t.Cancel()
}
