package task

// NOTE:
// there is a better way to handle tasks with interfaces

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"sync"
	"time"
)

var (
	taskMutex           = sync.RWMutex{}
	TaskDoesNotExistErr = errors.New("task does not exist")
	tasks               = make(map[string]*Task)
)

// DoesTaskExist checks if a task exists
func DoesTaskExist(id string) bool {
	taskMutex.RLock()
	defer taskMutex.RUnlock()
	_, ok := tasks[id]
	return ok
}

// CreateTask creates a task
func CreateTask(registeredSiteName, productInformation, size, profileId, proxyListId, accountId, taskType, taskMode string, delay int) string {
	taskMutex.Lock()
	defer taskMutex.Unlock()

	id := shortuuid.New()
	tInt := TypeConversionSTI(taskType)

	tasks[id] = &Task{
		ID:          id,
		TaskType:    tInt,
		Mode:        taskMode,
		Type:        registeredSiteName,
		Product:     productInformation,
		Size:        size,
		ProfileId:   profileId,
		ProxyListID: proxyListId,
		AccountId:   accountId,
		SiteId:      SiteConversionSTI[registeredSiteName],
		Delay:       time.Duration(delay) * time.Millisecond,
	}

	return id
}

// RemoveTask removes a task
func RemoveTask(id string) error {
	if !DoesTaskExist(id) {
		return TaskDoesNotExistErr
	}

	taskMutex.Lock()
	defer taskMutex.Unlock()

	// stop the task if active
	task := tasks[id]
	task.Cancel()

	delete(tasks, id)

	return nil
}

// GetTask gets a task
func GetTask(id string) (*Task, error) {
	if !DoesTaskExist(id) {
		return &Task{}, TaskDoesNotExistErr
	}

	taskMutex.RLock()
	defer taskMutex.RUnlock()

	return tasks[id], nil
}

func GetAllTasks() []*Task {
	taskMutex.RLock()
	defer taskMutex.RUnlock()

	tTasks := make([]*Task, 0)
	for _, task := range tasks {
		tTasks = append(tTasks, task)
	}

	return tTasks
}

func (t *Task) NotifyTasks(monitorData interface{}) error {
	for _, task := range GetAllTasks() {
		// if specific task is a checkout task
		// if requesting task is a monitor task
		// if both tasks are monitoring the same item
		// if both tasks are the same site
		if task.TaskType == 1 && t.TaskType == 0 && t.Product == task.Product && t.SiteId == task.SiteId {
			// works well on buzz, not athletesfoot?
			func() { // go func breaks?
				task.MonitorData <- monitorData
			}()
		}
	}

	return nil
}
