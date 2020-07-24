package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// Synchronizer is implemented by structs that do the actual synchronization actions.
type Synchronizer interface {
	GetConfig() *cfg.SynchConfig
	GetReporter() unifier.Reporter
	GetNodes() map[string]*node
	GetType() string
	IsSimulation() bool
	Run()
}
