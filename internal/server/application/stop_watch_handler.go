package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type stopWatchHandler struct {
	app *Application
}

func (h *stopWatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stop, ok := r.URL.Query()["stop"]
	if !ok || len(stop[0]) < 1 {
		log.Fatalln("[http request] ERROR: URL param 'stop' is missing.")
	}

	// A response channel can receive data of type 'error' or []byte.
	resChan := createResponseChannel()
	go h.app.stopWatch(resChan, stop[0])

	response := <-resChan
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshalling response.")
	}

	fmt.Fprintf(w, "%s", responseJSON)
}
