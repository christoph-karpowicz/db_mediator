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
	"github.com/christoph-karpowicz/unifier/internal/server/synch"
)

/*
Application is the main app object.
Contains all synchronization and database objects.
Starts a web server and handles all requests.
*/
type Application struct {
	dbs    db.Databases
	synchs synch.Synchs
}

// Init starts the application.
func (a *Application) Init() {
	a.dbs = make(db.Databases)
	a.dbs.Init()
	a.synchs = synch.CreateSynchs()
	a.synchs.Init()
	a.listen()
}

func (a *Application) listen() {
	http.Handle("/run", &runHandler{app: a})
	http.Handle("/stop", &stopHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

// synchronize carries out a synchronization requested by the client.
func (a *Application) synchronize(resChan chan interface{}, synchType string, synchKey string, simulation bool) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	synch, synchFound := a.synchs[synchKey]
	if !synchFound {
		panic("[synchronization search] '" + synchKey + "' not found.")
	}

	synch.SetSimulation()
	synch.SetReporter(report.CreateReport(synch))

	// Initialize synchronization.
	synch.Init(a.dbs)
	// Initialize report data structures.
	synch.GetReporter().Init()

	// Carry out all synch actions.
	synch.Run()
	synch.SetInitial(false)

	if synchType == "ongoing" {
		for synch.IsRunning() {
			fmt.Println("run")
			synch.Run()
			time.Sleep(5 * time.Second)
		}
	}

	// Gather and marshal results.
	synchReport, err := synch.GetReporter().Finalize()
	if err != nil {
		panic(err)
	}

	// Send the report to the http init handler.
	resChan <- synchReport

	synch.Reset()
}

// synchronize carries out a synchronization requested by the client.
func (a *Application) stop(resChan chan interface{}, synchKey string, all bool) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	if !all {
		for _, synch := range a.synchs {
			synch.Stop()
			synch.Reset()
		}
	} else {
		synch, synchFound := a.synchs[synchKey]
		if !synchFound {
			panic("[synchronization search] '" + synchKey + "' not found.")
		}

		synch.Stop()
		synch.Reset()
	}

	// Gather and marshal results.
	// synchReport, err := synch.GetReporter().Finalize()
	// if err != nil {
	// 	panic(err)
	// }

	// Send the report to the http init handler.
	// resChan <- synchReport
	resChan <- "stopped"
}

// synchronizeArray carries out aan array of synchronizations requested by the client.
// func (a *Application) synchronizeArray(synchType string, synchKeys []string, simulation bool) {
// 	for _, arg := range synchKeys {
// 		a.synchronize(synchType, arg, simulation)
// 	}
// }
