package logs

import "time"

func LogtailInitialize() {
	go func() {
		for {
			time.Sleep(time.Second * 15)
			flushLogs()
		}
	}()
}
