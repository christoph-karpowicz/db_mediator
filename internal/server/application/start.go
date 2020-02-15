package application

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var count int = 0

func tst(i int, w *http.ResponseWriter) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		fmt.Println(strconv.Itoa(i) + "*")
	}
	fmt.Fprintf(*w, "Done!")
}

type startHandler struct {
	app *Application
}

func (h *startHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

	h.app.Lang = "test222222"
	count++
	// go tst(count, &w)
	go h.app.Synchronize(synchType[0], synch[0], simulation)
}
