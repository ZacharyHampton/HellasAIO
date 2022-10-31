package quicktasking

import (
	"github.com/HellasAIO/HellasAIO/internal/loading"
	"github.com/HellasAIO/HellasAIO/internal/profile"
	"github.com/HellasAIO/HellasAIO/internal/proxy"
	"github.com/HellasAIO/HellasAIO/internal/task"
	taskmngr "github.com/HellasAIO/HellasAIO/internal/task_manager"
	"net/http"
	"net/url"
)

// only for fuel atm
func quicktaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/quicktask" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	msku := r.URL.Query().Get("product_id")
	if msku == "" {
		http.Error(w, "msku not provided.", http.StatusNotFound)
		return
	}

	siteId := r.URL.Query().Get("siteId")
	if siteId == "" {
		http.Error(w, "siteId not provided.", http.StatusNotFound)
		return
	}

	size, _ := url.QueryUnescape(r.URL.Query().Get("size"))
	if size == "" {
		http.Error(w, "size not provided.", http.StatusNotFound)
		return
	}

	proxyGroupId, _ := proxy.GetProxyGroupIDByName("main")
	var defaultAccount string
	var profileId string
	var err error

	switch siteId {
	case "1":
		defaultAccount = loading.Data.Settings.Settings.Quicktasks.Fuel.DefaultAccount
		profileId, err = profile.GetProfileIDByName(loading.Data.Settings.Settings.Quicktasks.Fuel.DefaultProfile)
		if err != nil {
			http.Error(w, "profile not found.", http.StatusNotFound)
			return
		}

		taskUUID := task.CreateTask(
			"fuelmonitor",
			msku,
			size,
			profileId,
			proxyGroupId,
			defaultAccount,
			"taskType",
			"fast",
			2000,
		)
		tObject, _ := task.GetTask(taskUUID)
		go taskmngr.RunTask(tObject)
	case "0":
		var taskMode = "msku"
		defaultAccount = loading.Data.Settings.Settings.Quicktasks.AthletesFoot.DefaultAccount
		profileId, err = profile.GetProfileIDByName(loading.Data.Settings.Settings.Quicktasks.AthletesFoot.DefaultProfile)
		if err != nil {
			http.Error(w, "profile not found.", http.StatusNotFound)
			return
		}

		monitorTaskUUID := task.CreateTask(
			"athletesfootmonitor",
			msku,
			size,
			profileId,
			proxyGroupId,
			defaultAccount,
			"monitor",
			taskMode,
			2000,
		)

		checkoutTaskUUID := task.CreateTask(
			"athletesfoot",
			msku,
			size,
			profileId,
			proxyGroupId,
			defaultAccount,
			"checkout",
			taskMode,
			2000,
		)

		monitorObject, _ := task.GetTask(monitorTaskUUID)
		checkoutObject, _ := task.GetTask(checkoutTaskUUID)
		go taskmngr.RunTask(monitorObject)
		go taskmngr.RunTask(checkoutObject)
	case "3":
		var taskMode = "msku"
		defaultAccount = loading.Data.Settings.Settings.Quicktasks.Buzzsneakers.DefaultAccount
		profileId, err = profile.GetProfileIDByName(loading.Data.Settings.Settings.Quicktasks.Buzzsneakers.DefaultProfile)
		if err != nil {
			http.Error(w, "profile not found.", http.StatusNotFound)
			return
		}

		monitorTaskUUID := task.CreateTask(
			"buzzsneakersmonitor",
			msku,
			size,
			profileId,
			proxyGroupId,
			defaultAccount,
			"monitor",
			taskMode,
			2000,
		)

		checkoutTaskUUID := task.CreateTask(
			"buzzsneakers",
			msku,
			size,
			profileId,
			proxyGroupId,
			defaultAccount,
			"checkout",
			taskMode,
			2000,
		)

		monitorObject, _ := task.GetTask(monitorTaskUUID)
		checkoutObject, _ := task.GetTask(checkoutTaskUUID)
		go taskmngr.RunTask(monitorObject)
		go taskmngr.RunTask(checkoutObject)
	}

	_, err = http.ResponseWriter(w).Write([]byte("Task created."))
	if err != nil {
		return
	}

}
