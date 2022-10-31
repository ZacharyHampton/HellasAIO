package quicktasking

import (
	"log"
	"net/http"
)

func start() {
	http.HandleFunc("/quicktask", quicktaskHandler)
	if err := http.ListenAndServe(":18638", nil); err != nil {
		log.Fatal(err)
	}
}
