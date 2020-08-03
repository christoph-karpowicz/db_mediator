package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Watchers is a collection of all watchers.
type Watchers map[string]*Watcher

// Init loads configs from files and validates them.
func (w *Watchers) Init() {
	w.getConfigs()
	w.validateConfigs()
}

func (w *Watchers) getConfigs() {
	var wtchCfgs []cfg.Config = cfg.GetWatcherConfigs()

	for i := 0; i < len(wtchCfgs); i++ {
		wtchCfg := wtchCfgs[i].(*cfg.WatcherConfig)
		(*w)[wtchCfg.Name] = &Watcher{cfg: wtchCfg, initial: true}
	}
}

// validateConfigs validates data imported from a config file.
func (w *Watchers) validateConfigs() {
	for _, wtch := range *w {
		(*wtch).GetConfig().Validate()
	}
}

// CreateWatchers constructor function for the Watchers struct.
func CreateWatchers() Watchers {
	synchs := make(Watchers)
	return synchs
}
