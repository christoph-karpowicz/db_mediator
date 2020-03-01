/*
Package application handles all initializations and
I/O of the app.
*/
package application

import (
	"net/http"

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
	dbs    *db.Databases
	synchs *synch.Synchs
}

// Init starts the application.
func (a *Application) Init() {
	a.dbs = &db.Databases{DBMap: make(map[string]*db.Database)}
	a.dbs.ImportYAML()
	a.dbs.ValidateYAML()
	a.synchs = synch.CreateSynchs()
	a.synchs.ImportYAMLDir()
	a.synchs.ValidateYAML()

	a.listen()
}

func (a *Application) listen() {
	http.Handle("/init", &initHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

// synchronize carries out a synchronization requested by the client.
func (a *Application) synchronize(resChan chan interface{}, synchType string, synchKey string, simulation bool) {
	defer func() {
		if r := recover(); r != nil {
			resChan <- r.(error)
		}
	}()

	// fmt.Printf("%s - %s\n", synchType, synchKey)
	synch, synchFound := a.synchs.SynchMap[synchKey]
	if !synchFound {
		panic("[synchronization search] '" + synchKey + "' not found.")
	}

	synch.Simulation = simulation
	synch.Rep = report.CreateReport(synch)

	// Initialize synchronization.
	synch.Init(a.dbs.DBMap)
	// Initialize report data structures.
	synch.Rep.Init()

	// Carry out all synch actions.
	synch.Synchronize()

	// Gather and marshal results.
	synchReport, err := synch.Rep.Finalize()
	if err != nil {
		panic("[report JSON marshalling] " + err.Error())
	}

	// Send the report to the http init handler.
	resChan <- synchReport
}

// synchronizeArray carries out aan array of synchronizations requested by the client.
// func (a *Application) synchronizeArray(synchType string, synchKeys []string, simulation bool) {
// 	for _, arg := range synchKeys {
// 		a.synchronize(synchType, arg, simulation)
// 	}
// }
