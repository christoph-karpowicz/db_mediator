package sim

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/synch"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// PairSimulator is implemented by the pair struct from the synch package
// and passed to Add* methods.
type PairSimulator interface {
	SimIdleString() string
	SimInsertString() string
	SimUpdateString() string
}

// Simulation is basically a report about what will happen after an actual synchronization is launched.
type Simulation struct {
	idle    []string
	inserts []string
	updates []string
}

// CreateSimulation creates a Simulation instance.
func CreateSimulation(s *synch.Synch) unifier.Simulator {
	var newSimulation Simulation = Simulation{}

	return &newSimulation
}

// AddIdle adds a single idle to the Simulation.
// Idle means two records that have been paired, but no action will be carried out because the relevant data is the same.
func (s *Simulation) AddIdle(pair unifier.Synchronizer) (bool, error) {
	var str string = pair.(synch.Pair).SimIdleString()
	s.idle = append(s.idle, str)
	fmt.Print(str)
	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert(pair unifier.Synchronizer) (bool, error) {
	var str string = pair.(synch.Pair).SimInsertString()
	s.inserts = append(s.inserts, str)
	fmt.Print(str)
	return false, nil
}

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(pair unifier.Synchronizer) (bool, error) {
	var str string = pair.(synch.Pair).SimUpdateString()
	s.updates = append(s.updates, str)
	fmt.Print(str)
	return false, nil
}
