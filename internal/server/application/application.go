/*
Package application handles all initializations and
I/O of the app.
*/
package application

import (
	"fmt"
	"net/http"
	"strings"
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
	synch.GetHistory().Init()

	// Carry out all synch actions.
	if !simulation && synchType == "ongoing" {
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

func (a *Application) runOngoing(synch *synch.Synch) {
	for synch.IsInitial() || synch.IsRunning() {
		fmt.Println("run")
		synch.Run()
		synch.SetInitial(false)
		time.Sleep(1 * time.Second)
	}
}

// synchronize carries out a synchronization requested by the client.
func (a *Application) stop(resChan chan interface{}, synchKey string, all bool) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	var response string
	synchsStopped := make([]string, 0)
	if !all {
		for skey, synch := range a.synchs {
			if synch.IsRunning() {
				synch.Stop()
				synch.Reset()
				synchsStopped = append(synchsStopped, skey)
			}
		}
	} else {
		synch, synchFound := a.synchs[synchKey]
		if !synchFound {
			panic("[synchronization search] '" + synchKey + "' not found.")
		}

		if synch.IsRunning() {
			synch.Stop()
			synch.Reset()
			synchsStopped = append(synchsStopped, synchKey)
		}
	}

	if len(synchsStopped) > 0 {
		synchWord := "Synch"
		if len(synchsStopped) > 1 {
			synchWord += "s"
		}

		for _, synchStopped := range synchsStopped {
			// Gather and marshal results.
			synchReport, err := a.synchs[synchStopped].GetHistory().GenerateReport()
			if err != nil {
				panic(err)
			}
			response += string(synchReport)
		}

		response += fmt.Sprintf("\n\n%s %s stopped.", synchWord, strings.Join(synchsStopped, ", "))
	} else if all && len(synchsStopped) == 0 {
		response = "No running synchs found."
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
