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

type StartHandler struct {
	app *Application
}

func (h *StartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	synchs, ok := r.URL.Query()["synch"]

	if !ok || len(synchs[0]) < 1 {
		log.Println("Url Param 'synch' is missing")
		return
	}

	// // synch := keys[0]
	// log.Println(r.URL.Query())

	// log.Println("Url Param 'synch' is: " + string(synch))

	// fmt.Fprintf(w, "Welcome to my website!")
	fmt.Println("h.app.Lang::::::")
	fmt.Println(h.app.Lang)
	h.app.Lang = "test222222"
	count++
	// go tst(count, &w)
	go h.app.synchronizeArray(synchs)
}
