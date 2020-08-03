package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Synchs is a collection of all synchronizations.
type Synchs map[string]*Synch

// Init loads configs from files and validates them.
func (s *Synchs) Init() {
	s.getConfigs()
	s.validateConfigs()
}

func (s *Synchs) getConfigs() {
	var synchCfgs []cfg.Config = cfg.GetSynchConfigs()

	for i := 0; i < len(synchCfgs); i++ {
		synchCfg := synchCfgs[i].(*cfg.SynchConfig)
		(*s)[synchCfg.Name] = &Synch{cfg: synchCfg, initial: true}
	}
}

// validateConfigs validates data imported from a config file.
func (s *Synchs) validateConfigs() {
	for _, synch := range *s {
		(*synch).GetConfig().Validate()
	}
}

// CreateSynchs constructor function for the Synchs struct.
func CreateSynchs() Synchs {
	synchs := make(Synchs)
	return synchs
}
