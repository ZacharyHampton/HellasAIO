package logs

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/loading"
	"github.com/HellasAIO/HellasAIO/internal/task"
)

func Log(t *task.Task, data ...interface{}) {
	siteName := t.Type
	taskType := task.TypeConversionITS(t.TaskType)
	taskMode := t.Mode
	taskProduct := t.Product
	stringData := fmt.Sprint(data...)
	authKey := loading.Data.Settings.Settings.AuthKey

	go LogLogTail(siteName, taskType, taskMode, taskProduct, stringData, authKey)
	fmt.Println(fmt.Sprintf("[%s] [%s] [%s] %s: %v", siteName, taskType, taskMode, taskProduct, stringData))
}
