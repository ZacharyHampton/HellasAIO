package logs

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/internal/version"
	"github.com/HellasAIO/HellasAIO/third_party/hclient"
	"github.com/getsentry/sentry-go"
	"sync"
	"time"
)

var (
	Logs      = make(map[LogtailData]int)
	logsMutex = sync.RWMutex{}
)

func LogLogTail(siteName, taskType, taskMode, taskProduct, message, key string) {
	defer func() {
		err := recover()

		if err != nil {
			localHub := sentry.CurrentHub().Clone()
			localHub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("logtail", "log")
			})

			localHub.Recover(err)
			sentry.Flush(time.Second * 5)
		}
	}()

	data := LogtailData{
		AuthKey:     key,
		SiteName:    siteName,
		TaskType:    taskType,
		TaskMode:    taskMode,
		TaskProduct: taskProduct,
		Version:     version.Version,
		Message:     message,
	}

	logsMutex.RLock()
	_, logExists := Logs[data]
	logsMutex.RUnlock()

	logsMutex.Lock()
	defer logsMutex.Unlock()
	if logExists {
		Logs[data]++
	} else {
		Logs[data] = 1
	}
}

func getRequestBodyForFlush() []LogtailData {
	var requestBody []LogtailData
	for key, value := range Logs {
		key.Count = value
		requestBody = append(requestBody, key)
	}

	return requestBody
}

func flushLogs() {
	defer func() {
		err := recover()

		if err != nil {
			localHub := sentry.CurrentHub().Clone()
			localHub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("logtail", "flush")
			})

			localHub.Recover(err)
			sentry.Flush(time.Second * 5)
		}
	}()

	logsMutex.RLock()

	if len(Logs) > 0 {
		requestBody := getRequestBodyForFlush()
		logsMutex.RUnlock()

		client, err := hclient.NewClient()
		if err != nil {
			return
		}

		_, err = client.NewRequest().
			SetMethod("POST").
			SetURL("https://in.logtail.com").
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "*/*").
			SetHeader("Authorization", "Bearer authkey").
			SetJSONBody(requestBody).
			Do()

		if err != nil {
			return
		} else {
			logsMutex.Lock()
			Logs = make(map[LogtailData]int)
			logsMutex.Unlock()
		}

		if utils.Debug {
			fmt.Println("Logs status code:", client.LatestResponse.StatusCode())
		}
	} else {
		logsMutex.RUnlock()
	}
}
