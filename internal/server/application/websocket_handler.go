package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type webSocketHandler struct {
	app *Application
}

func (wsh *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Client Connected")
	log.Println(string(wsh.app.listWatchersToJSON()))

	if wsh.app == nil {
		log.Println("err")
	}
	err = ws.WriteMessage(1, wsh.app.listWatchersToJSON())
	if err != nil {
		log.Println(err)
	}

	go wsh.wsReader(ws)
}

func (wsh *webSocketHandler) wsReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
