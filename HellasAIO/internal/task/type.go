package task

import (
	"errors"
	"reflect"
)

var (
	DoneTaskState  TaskState = "done"
	ErrorTaskState TaskState = "error"

	TaskTypeDoesNotExistErr    = errors.New("task type does not exist")
	TaskHandlerDoesNotExistErr = errors.New("task handler does not exist")

	taskTypes         = make(map[string]*TaskType)
	SiteConversionSTI = make(map[string]int)
	SiteConversionITS = make(map[int][]string)
)

func TypeConversionSTI(taskType string) int {
	if taskType == "checkout" {
		return 1
	}

	if taskType == "monitor" {
		return 0
	}

	return -1
}

func TypeConversionITS(taskType int) string {
	if taskType == 1 {
		return "checkout"
	}

	if taskType == 0 {
		return "monitor"
	}

	return ""
}

// taskType ? 0:1 | 0 = Monitor | 1 = Checkout

// RegisterTaskType register task type
func RegisterTaskType(registerSiteName string, siteId int) *TaskType {
	taskTypes[registerSiteName] = &TaskType{
		handlers: make(TaskReflectMap),
	}

	SiteConversionSTI[registerSiteName] = siteId
	SiteConversionITS[siteId] = append(SiteConversionITS[siteId], registerSiteName)

	return taskTypes[registerSiteName]
}

// DoesTaskTypeExist check if task type exists
func DoesTaskTypeExist(taskType string) bool {
	_, ok := taskTypes[taskType]
	return ok
}

// GetTaskType gets a task type
func GetTaskType(taskType string) (*TaskType, error) {
	if !DoesTaskTypeExist(taskType) {
		return &TaskType{}, TaskTypeDoesNotExistErr
	}
	return taskTypes[taskType], nil
}

// HasHandlers check if there are handlers
func (t *TaskType) HasHandlers() bool {
	return len(t.handlers) > 0
}

func (t *TaskType) addHandler(handlerName TaskState, handler interface{}) {
	t.handlers[string(handlerName)] = reflect.ValueOf(handler)
}

// AddHandlers adds multiple handles to a task type
func (t *TaskType) AddHandlers(handlers TaskHandlerMap) {
	for handlerName, handler := range handlers {
		if t.internalType == nil {
			handleTypes := reflect.TypeOf(handler)
			// func (t *task.Task, internal *SiteInternal) task.TaskState

			// we want the second one because the first one (0 index) will be task.Task type
			handleType := handleTypes.In(1)

			t.internalType = handleType
		}

		t.addHandler(handlerName, handler)
	}
}

// GetHandler gets a handler by handler name
func (t *TaskType) GetHandler(handlerName TaskState) (reflect.Value, error) {
	handler, ok := t.handlers[string(handlerName)]

	if !ok {
		return reflect.Value{}, TaskHandlerDoesNotExistErr
	}

	return handler, nil
}

// GetFirstHandlerState gets the first handler state
func (t *TaskType) GetFirstHandlerState() TaskState {
	return t.firstHandlerState
}

// SetFirstHandlerState sets the first handler state
func (t *TaskType) SetFirstHandlerState(firstHandlerState TaskState) {
	t.firstHandlerState = firstHandlerState
}

// GetInternalType gets the internal type
func (t *TaskType) GetInternalType() reflect.Type {
	return t.internalType
}
