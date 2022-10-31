package main

import (
	"fmt"
	"github.com/HellasAIO/HellasAIO/internal/auth"
	"github.com/HellasAIO/HellasAIO/internal/loading"
	"github.com/HellasAIO/HellasAIO/internal/logs"
	"github.com/HellasAIO/HellasAIO/internal/rich_presence"
	sntry "github.com/HellasAIO/HellasAIO/internal/sentry"
	"github.com/HellasAIO/HellasAIO/internal/sessions"
	"github.com/HellasAIO/HellasAIO/internal/site_parsing"
	"github.com/HellasAIO/HellasAIO/internal/task"
	taskmngr "github.com/HellasAIO/HellasAIO/internal/task_manager"
	"github.com/HellasAIO/HellasAIO/internal/title"
	"github.com/HellasAIO/HellasAIO/internal/ws_quicktasking"
	athletesfootmonitor "github.com/HellasAIO/HellasAIO/monitors/athletesfoot"
	buzzsneakersmonitor "github.com/HellasAIO/HellasAIO/monitors/buzzsneakers"
	fuelmonitor "github.com/HellasAIO/HellasAIO/monitors/fuel"
	"github.com/HellasAIO/HellasAIO/sites/athletesfoot"
	"github.com/HellasAIO/HellasAIO/sites/buzzsneakers"
	"github.com/getsentry/sentry-go"
	"os"
	"time"
)

func main() {
	athletesfootmonitor.Initialize()
	athletesfoot.Initialize()
	fuelmonitor.Initialize()
	buzzsneakersmonitor.Initialize()
	buzzsneakers.Initialize()

	loading.Initialize()
	auth.Initialize()
	sntry.Initialize()
	ws_quicktasking.Initialize()
	sessions.Initialize()
	title.Initialize()
	rich_presence.Initialize()
	logs.LogtailInitialize()

	defer sentry.Recover()
	defer sentry.Flush(2 * time.Second)

	fmt.Printf(`1. AthletesFoot
2. Fuel
3. Slamdunk
4. Buzzsneakers
5. Europe Sports
9. Exit` + "\n\n\n")

	for {
		var input string
		fmt.Scanln(&input)

		if input == "9" {
			os.Exit(0)
		}

		data := site_parsing.Parse(input)
		if data == nil {
			continue
		}

		rich_presence.SetSite(data.SiteID)
		for _, taskUUID := range loading.Data.Tasks.Tasks[data.SiteID] {
			taskObject, err := task.GetTask(taskUUID)
			if err != nil {
				fmt.Println("Failed to get task: ", err.Error())
				continue
			}

			if data.TaskType != 2 {
				if taskObject.TaskType != data.TaskType {
					continue
				}
			}

			if data.Action == 0 {
				if !taskObject.Active {
					go taskmngr.RunTask(taskObject)
				}
			} else if data.Action == 1 {
				if taskObject.Active {
					taskmngr.StopTask(taskObject)
				}
			}
		}
	}
}
