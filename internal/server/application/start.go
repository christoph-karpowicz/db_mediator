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
		log.Println("Url Param 'synchType' is missing")
		return
	}

	synch, ok := r.URL.Query()["synch"]
	if !ok || len(synch[0]) < 1 {
		log.Println("Url Param 'synch' is missing")
		return
	}

	h.app.Lang = "test222222"
	count++
	// go tst(count, &w)
	go h.app.synchronize(synchType[0], synch[0])
}
