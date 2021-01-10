/*
Package application handles all initializations and
I/O of the app.
*/
package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	synchPkg "github.com/christoph-karpowicz/unifier/internal/server/synch"
)

/*
Application is the main app object.
Contains all synchronization and database objects.
Starts a web server and handles all requests.
*/
type Application struct {
	dbs    db.Databases
	synchs synchPkg.Synchs
}

// Init starts the application.
func (a *Application) Init() {
	a.dbs = make(db.Databases)
	a.dbs.Init()
	a.synchs = synchPkg.CreateSynchs()
	a.synchs.Init()
	a.listen()
}

func (a *Application) listen() {
	http.Handle("/", &frontHandler{app: a})
	http.Handle("/ws/", &webSocketHandler{app: a})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/build/static"))))
	http.Handle("/runSynch", &runSynchHandler{app: a})
	http.Handle("/stopSynch", &stopSynchHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

// runSynch carries out a synchronization run requested by the client.
func (a *Application) runSynch(responseChan chan *response, synchType string, synchName string, isSimulation bool) {
	defer func() {
		if r := recover(); r != nil {
			responseChan <- createResponse(r.(error))
		}
	}()

	synch, synchFound := a.synchs[synchName]
	if !synchFound {
		panic("[synchronization search] '" + synchName + "' not found.")
	}

	synch.SetSimulation(isSimulation)

	// Initialize synchronization.
	synchID := synch.Init(a.dbs, synchType)

	// Carry out all synch actions.
	if !isSimulation && synch.GetType() == synchPkg.ONGOING {
		go a.runSynchLoop(synch)
		responseChan <- createResponse(fmt.Sprintf("Synch %s started with ID %s.", synchName, synchID))
	} else {
		synch.Run()
		synchResponse := synch.Flush()
		responseChan <- createResponse(synchResponse)
		synch.Reset()
	}
}

func (a *Application) runSynchLoop(synch *synchPkg.Synch) {
	for synch.IsInitial() || synch.IsRunning() {
		fmt.Println("run synch")
		synch.Run()
		synch.SetInitial(false)
		time.Sleep(1 * time.Second)
	}
}

// stopSynch stops a specified synchronization.
func (a *Application) stopSynch(responseChan chan *response, synchName string) {
	defer func() {
		if r := recover(); r != nil {
			responseChan <- createResponse(r.(error))
		}
	}()

	var synchResponse interface{}
	synch, synchFound := a.synchs[synchName]
	if !synchFound {
		synchResponse = fmt.Sprintf("[synchronization search] \"%s\" not found.", synchName)
	} else if synch.IsRunning() {
		synch.Stop()
		synchResponse = synch.Flush()
		synch.Reset()
	} else {
		synchResponse = fmt.Sprintf("Synch \"%s\" is not running.", synchName)
	}

	fmt.Println(synchResponse)
	// Send the reponse to the http init handler.
	responseChan <- createResponse(synchResponse)
}

func (a *Application) listSynchs() []string {
	synchList := make([]string, 0)
	for name := range a.synchs {
		synchList = append(synchList, name)
	}
	return synchList
}

func (a *Application) listSynchsToJSON() []byte {
	synchList := a.listSynchs()
	synchListJSON, err := json.Marshal(synchList)
	if err != nil {
		panic(err)
	}
	return synchListJSON
}
