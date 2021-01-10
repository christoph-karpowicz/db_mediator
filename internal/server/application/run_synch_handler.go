package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type runSynchHandler struct {
	app *Application
}

func (h *runSynchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	synchType, ok := r.URL.Query()["type"]
	if !ok || len(synchType[0]) < 1 {
		log.Fatalln("[http request] ERROR: URL param 'type' is missing.")
	}

	run, ok := r.URL.Query()["run"]
	if !ok || len(run[0]) < 1 {
		log.Fatalln("[http request] ERROR: URL param 'run' is missing.")
	}

	simulationStr, ok := r.URL.Query()["simulation"]
	if !ok {
		simulationStr = []string{"false"}
	}
	simulation, err := strconv.ParseBool(simulationStr[0])
	if err != nil {
		log.Fatalln("[http request] ERROR: Wrong 'simulation' URL param value.")
	}

	if simulation && synchType[0] == "ongoing" {
		log.Fatalln("[http request] ERROR: Cannot start an ongoing synchronization simulation.")
	}

	resChan := createResponseChannel()
	go h.app.runSynch(resChan, synchType[0], run[0], simulation)

	response := <-resChan
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic("Error while marshalling response.")
	}

	fmt.Fprintf(w, "%s", responseJSON)
}
