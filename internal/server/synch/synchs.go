package synch

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Synchs is a collection of all synchronizations.
type Synchs struct {
	synchCfgs []*cfg.SynchConfig
	SynchMap  map[string]*Synch
}

func (s *Synchs) Init() {
	s.getConfigs()
	s.validateConfigs()
	s.assignConfigs()
}

func (s *Synchs) getConfigs() {
	var synchConfigArray []*cfg.SynchConfig = cfg.GetSynchConfigs()
	s.synchCfgs = synchConfigArray
}

func (s *Synchs) assignConfigs() {
	for i := 0; i < len(s.synchCfgs); i++ {
		synchCfg := s.synchCfgs[i]
		s.SynchMap[s.synchCfgs[i].Name] = &Synch{Cfg: synchCfg, initial: true}
		// fmt.Printf("val: %s\n", dbDataArr.Databases[i].Name)
	}
}

// validateConfigs validates data imported from a config file.
func (s *Synchs) validateConfigs() {
	fmt.Println("Synch YAML file validation...")
	for _, synch := range s.SynchMap {
		(*synch).GetConfig().Validate()
	}
	fmt.Println("... passed.")
}

// CreateSynchs constructor function for the Synchs struct.
func CreateSynchs() *Synchs {
	synchs := &Synchs{SynchMap: make(map[string]*Synch)}
	return synchs
}
