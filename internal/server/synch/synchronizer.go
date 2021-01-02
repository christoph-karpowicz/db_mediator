package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

// Synchronizer is implemented by structs that do the actual synchronization actions.
type Synchronizer interface {
	GetConfig() *cfg.SynchConfig
	GetIteration() *iteration
	GetNodes() map[string]*node
	GetType() synchType
	IsSimulation() bool
	Run()
}
