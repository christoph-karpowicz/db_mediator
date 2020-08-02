/*
Package synch handles all data sychronization.
*/
package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

// Watcher represents an individual watcher configration.
// It holds all configuration from an .yaml file, raw and parsed.
type Watcher struct {
	cfg     *cfg.WatcherConfig
	dbs     map[string]*db.Database
	tables  map[string]*table
	nodes   map[string]*node
	running bool
	initial bool
	History *History
}

// GetConfig returns the synch config struct.
func (w *Watcher) GetConfig() *cfg.WatcherConfig {
	return w.cfg
}
