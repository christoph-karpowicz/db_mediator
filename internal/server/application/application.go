/*
Package application handles all initializations and
I/O of the app.
*/
package application

import (
	"fmt"
	"net/http"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/synch"
)

/*
Application is the main app object.
Contains all synchronization and database objects.
Starts a web server and handles all requests.
*/
type Application struct {
	Lang   string
	dbs    *db.Databases
	synchs *synch.Synchs
}

// Init starts the application.
func (a *Application) Init() {
	a.dbs = &db.Databases{DBMap: make(map[string]*db.Database)}
	a.dbs.ImportJSON()
	a.dbs.ValidateJSON()
	a.synchs = synch.CreateSynchs()
	a.synchs.ImportJSONDir()
	a.synchs.ValidateJSON()

	a.listen()
}

func (a *Application) listen() {
	http.Handle("/init", &startHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

func (a *Application) synchronize(synchType string, synchKey string) {
	fmt.Printf("%s - %s\n", synchType, synchKey)
	synch := a.synchs.SynchMap[synchKey]
	if synch == nil {
		panic("Synch '" + synchKey + "' not found.")
	}

	synch.Init(a.dbs.DBMap)
	synch.SynchPairs()

	// fmt.Println(*synch)
}

func (a *Application) synchronizeArray(synchType string, synchKeys []string) {
	for _, arg := range synchKeys {
		a.synchronize(synchType, arg)
	}
}
