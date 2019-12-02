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
	a.synchs.ImportJSON()
	a.synchs.ValidateJSON()

	a.Lang = "sdsadsdsad"
	a.listen()
}

func (a *Application) listen() {
	http.Handle("/start", &StartHandler{app: a})
	http.ListenAndServe(":8000", nil)
}

func (a *Application) synchronize(synchKey string) {
	synch := a.synchs.SynchMap[synchKey]
	if synch == nil {
		panic("Synch '" + synchKey + "' not found.")
	}

	synch.SetDatabases(a.dbs.DBMap)

	fmt.Println(*synch)
}

func (a *Application) synchronizeArray(synchKeys []string) {
	for _, arg := range synchKeys {
		a.synchronize(arg)
	}
}
