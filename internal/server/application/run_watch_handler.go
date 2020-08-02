package application

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type runWatchHandler struct {
	app *Application
}

func (h *runWatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	run, ok := r.URL.Query()["run"]
	if !ok || len(run[0]) < 1 {
		run = []string{""}
		// log.Fatalln("[http request] ERROR: URL param 'run' is missing.")
	}

	// A response channel can receive data of type 'error' or []byte.
	resChan := make(chan interface{})
	go h.app.runWatch(resChan, run[0])

	response := createResponse(<-resChan)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshalling response.")
	}

	fmt.Fprintf(w, "%s", responseJSON)
}
