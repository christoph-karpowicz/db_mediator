package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type stopHandler struct {
	app *Application
}

func (h *stopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stop, ok := r.URL.Query()["stop"]
	if !ok || len(stop[0]) < 1 {
		stop = []string{""}
		// log.Fatalln("[http request] ERROR: URL param 'stop' is missing.")
	}

	allStr, ok := r.URL.Query()["all"]
	if !ok {
		allStr = []string{"false"}
	}
	all, err := strconv.ParseBool(allStr[0])
	if err != nil {
		log.Fatalln("[http request] ERROR: Wrong 'all' URL param value.")
	}

	// A response channel can receive data of type 'error' or []byte.
	resChan := make(chan interface{})
	go h.app.stop(resChan, stop[0], all)

	response := createResponse(<-resChan)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshalling response.")
	}

	fmt.Fprintf(w, "%s", responseJSON)
}
