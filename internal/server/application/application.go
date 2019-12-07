package application

import (
	"fmt"
	"net/http"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/synch"
)

type Application struct {
	Lang   string
	dbs    *db.Databases
	synchs *synch.Synchs
}

func (a *Application) Init() {
	a.dbs = &db.Databases{DBMap: make(map[string]*db.Database)}
	a.dbs.ImportJSON()
	a.dbs.ValidateJSON()
	a.synchs = &synch.Synchs{SynchMap: make(map[string]*synch.Synch)}
	a.synchs.ImportJSONDir()
	a.synchs.ValidateJSON()

	a.Lang = "sdsadsdsad"
	a.listen()
}

func (a *Application) listen() {
	http.Handle("/init", &StartHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

func (a *Application) synchronize(synchType string, synchKey string) {
	fmt.Printf("%s - %s\n", synchType, synchKey)
	synch := a.synchs.SynchMap[synchKey]
	if synch == nil {
		panic("Synch '" + synchKey + "' not found.")
	}

	synch.SetDatabases(a.dbs.DBMap)
	synch.SelectData()

	fmt.Println(*synch)
}

func (a *Application) synchronizeArray(synchType string, synchKeys []string) {
	for _, arg := range synchKeys {
		a.synchronize(synchType, arg)
	}
}
