package synch

import (
	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
)

// Synchronizer is implemented by structs that do the actual synchronization actions.
type Synchronizer interface {
	GetConfig() *cfg.SynchConfig
	GetIteration() *iteration
	GetNodes() map[string]*node
	GetRawMappings() []map[string]string
	GetType() synchType
	IsSimulation() bool
	Run()
}
