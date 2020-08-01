/*
Package application handles all initializations and
I/O of the app.
*/
package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/report"
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("front/build/static"))))
	http.Handle("/runSynch", &runHandler{app: a})
	http.Handle("/stopSynch", &stopHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

// run carries out a synchronization run requested by the client.
func (a *Application) run(resChan chan interface{}, synchType string, synchKey string, simulation bool) {
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
	synch.GetHistory().SetReporter(report.CreateReporter(synch))
	synch.GetHistory().Init(synch)

	// Carry out all synch actions.
	if !simulation && synch.GetType() == synchPkg.ONGOING {
		go a.runOngoing(synch)
		resChan <- fmt.Sprintf("Synch %s started.", synchKey)
	} else {
		synch.Run()

		// Gather and marshal results.
		synchReport, err := synch.GetHistory().GenerateReport()
		if err != nil {
			panic(err)
		}

		resChan <- synchReport
		synch.Reset()
	}
}

func (a *Application) runOngoing(synch *synchPkg.Synch) {
	for synch.IsInitial() || synch.IsRunning() {
		fmt.Println("run")
		synch.Run()
		synch.SetInitial(false)
		time.Sleep(1 * time.Second)
	}
}

// synchronize carries out a synchronization requested by the client.
func (a *Application) stop(resChan chan interface{}, synchKey string) {
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
		synchReport, err := a.synchs[synchKey].GetHistory().GenerateReport()
		if err != nil {
			panic(err)
		}

		synch.Stop()
		synch.Reset()

		response = synchReport
	} else {
		response = fmt.Sprintf("Synch %s is not running.", synchKey)
	}

	// Gather and marshal results.
	// synchReport, err := synch.GetHistory().GenerateReport()
	// if err != nil {
	// 	panic(err)
	// }

	// Send the reponse to the http init handler.
	resChan <- response
}

// synchronizeArray carries out aan array of synchronizations requested by the client.
// func (a *Application) synchronizeArray(synchType string, synchKeys []string, simulation bool) {
// 	for _, arg := range synchKeys {
// 		a.synchronize(synchType, arg, simulation)
// 	}
// }
