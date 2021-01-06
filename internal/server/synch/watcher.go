/*
Package synch handles all data sychronization.
*/
package synch

import (
	"fmt"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

// Watcher represents an individual watcher configration.
// It holds all configuration from an .yaml file, raw and parsed.
type Watcher struct {
	cfg     *cfg.WatcherConfig
	dbStore *dbStore
	running bool
	initial bool
}

// Init prepares the watcher by fetching all necessary data
// and parsing it.
func (w *Watcher) Init(DBMap map[string]*db.Database) {
	tStart := time.Now()
	w.dbStore = &dbStore{}
	w.dbStore.Init(DBMap, w.cfg.Nodes)
	fmt.Println("Watcher init finished in: ", time.Since(tStart).String())
}

// GetConfig returns the synch config struct.
func (w *Watcher) GetConfig() *cfg.WatcherConfig {
	return w.cfg
}

func (w *Watcher) IsInitial() bool {
	return w.initial
}

// SetInitial sets the initial struct field indicating whether
// it'w the first run of the synch.
func (w *Watcher) SetInitial(ini bool) {
	w.initial = ini
}

func (w *Watcher) IsRunning() bool {
	return w.running
}

// Run executes a single run of the watcher.
func (w *Watcher) Run() {
	w.running = true

	// s.selectData()
	// s.pairData()
	// s.synchronize()
	// s.flush()
}
