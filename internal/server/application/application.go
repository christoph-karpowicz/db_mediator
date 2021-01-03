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
	dbs      db.Databases
	synchs   synchPkg.Synchs
	watchers synchPkg.Watchers
}

// Init starts the application.
func (a *Application) Init() {
	a.dbs = make(db.Databases)
	a.dbs.Init()
	a.synchs = synchPkg.CreateSynchs()
	a.synchs.Init()
	a.watchers = synchPkg.CreateWatchers()
	a.watchers.Init()
	a.listen()
}

func (a *Application) listen() {
	http.Handle("/", &frontHandler{app: a})
	http.Handle("/ws/", &webSocketHandler{app: a})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/build/static"))))
	http.Handle("/runSynch", &runSynchHandler{app: a})
	http.Handle("/stopSynch", &stopSynchHandler{app: a})
	http.Handle("/runWatch", &runWatchHandler{app: a})
	http.Handle("/stopWatch", &stopWatchHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

// runSynch carries out a synchronization run requested by the client.
func (a *Application) runSynch(resChan chan interface{}, synchType string, synchKey string, simulation bool) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		resChan <- r.(error)
	// 	}
	// }()

	synch, synchFound := a.synchs[synchKey]
	if !synchFound {
		panic("[synchronization search] '" + synchKey + "' not found.")
	}

	synch.SetSimulation(simulation)

	// Initialize synchronization.
	synch.Init(a.dbs, synchType)
	// Initialize history data structures.
	// synch.GetHistory().SetReporter(report.CreateReporter(synch))
	// synch.GetHistory().Init(synch)

	// Carry out all synch actions.
	if !simulation && synch.GetType() == synchPkg.ONGOING {
		go a.runSynchLoop(synch)
		resChan <- fmt.Sprintf("Synch %s started.", synchKey)
	} else {
		synchResponse := synch.Run()
		resChan <- synchResponse
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
func (a *Application) stopSynch(resChan chan interface{}, synchKey string) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	var response interface{}
	synch, synchFound := a.synchs[synchKey]
	if !synchFound {
		panic("[synchronization search] '" + synchKey + "' not found.")
	}

	if synch.IsRunning() {
		// Gather and marshal results.
		// synchReport, err := synch.GetHistory().GenerateReport()
		// if err != nil {
		// 	panic(err)
		// }

		synch.Stop()
		synch.Reset()

		response = "synchReport"
		// response = synchReport
	} else {
		response = fmt.Sprintf("Synch %s is not running.", synchKey)
	}

	// Send the reponse to the http init handler.
	resChan <- response
}

// runWatch starts a watcher.
func (a *Application) runWatch(resChan chan interface{}, watcherKey string) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		resChan <- r.(error)
	// 	}
	// }()

	var response interface{}
	watcher, watcherFound := a.watchers[watcherKey]
	if !watcherFound {
		panic("[watcher search] '" + watcherKey + "' not found.")
	}

	if watcher.IsRunning() {
		response = fmt.Sprintf("Watcher %s is already running.", watcherKey)
	} else {
		watcher.Init(a.dbs)
		go a.runWatchLoop(watcher)
		resChan <- fmt.Sprintf("Watcher %s started.", watcherKey)
	}

	// Send the reponse to the http init handler.
	resChan <- response
}

func (a *Application) runWatchLoop(watcher *synchPkg.Watcher) {
	for watcher.IsInitial() || watcher.IsRunning() {
		fmt.Println("run watch")
		watcher.Run()
		watcher.SetInitial(false)
		time.Sleep(1 * time.Second)
	}
}

// stopWatch stops a specified watcher.
func (a *Application) stopWatch(resChan chan interface{}, watcherKey string) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	var response interface{}
	watcher, watcherFound := a.watchers[watcherKey]
	if !watcherFound {
		panic("[watcher search] '" + watcherKey + "' not found.")
	}

	if watcher.IsRunning() {
		// watcher.Stop()
		// watcher.Reset()
		response = fmt.Sprintf("Watch %s stopped.", watcherKey)
	} else {
		response = fmt.Sprintf("Watch %s is not running.", watcherKey)
	}

	// Send the reponse to the http init handler.
	resChan <- response
}

func (a *Application) listWatchers() []string {
	watcherList := make([]string, 0)
	for name := range a.watchers {
		watcherList = append(watcherList, name)
	}
	return watcherList
}

func (a *Application) listWatchersToJSON() []byte {
	watcherList := a.listWatchers()
	watcherListJSON, err := json.Marshal(watcherList)
	if err != nil {
		panic(err)
	}
	return watcherListJSON
}
