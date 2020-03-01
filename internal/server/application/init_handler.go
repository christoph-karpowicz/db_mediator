package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type initHandler struct {
	app *Application
}

func (h *initHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	synchType, ok := r.URL.Query()["type"]
	if !ok || len(synchType[0]) < 1 {
		log.Fatalln("[http request] ERROR: URL param 'type' is missing.")
	}

	synch, ok := r.URL.Query()["synch"]
	if !ok || len(synch[0]) < 1 {
		log.Fatalln("[http request] ERROR: URL param 'synch' is missing.")
	}

	simulationStr, ok := r.URL.Query()["simulation"]
	if !ok {
		simulationStr = []string{"false"}
	}
	simulation, err := strconv.ParseBool(simulationStr[0])
	if err != nil {
		log.Fatalln("[http request] ERROR: Wrong 'simulation' URL param value.")
	}

	// A reponse channel can receive data of type 'error' or []byte.
	resChan := make(chan interface{})
	go h.app.synchronize(resChan, synchType[0], synch[0], simulation)

	response := createResponse(<-resChan)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshalling response.")
	}

	fmt.Fprintf(w, "%s", responseJSON)
}
