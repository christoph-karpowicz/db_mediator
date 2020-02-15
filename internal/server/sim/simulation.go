package sim

import "fmt"

// Simulator enables creating strings which represent data flows
// needed for the simulation directly in the synch package (avoiding a package cycle).
type Simulator interface {
	SimString() string
}

// PairSimulator is implemented by the pair struct from the synch package
// and passed to Add* methods.
type PairSimulator interface {
	SimIdleString() string
	SimInsertString() string
	SimUpdateString() string
}

// Simulation is basically a report about what will happen after an actual synchronization is launched.
type Simulation struct {
	synchName string
	updates   []string
	inserts   []string
}

// CreateSimulation creates a Simulation instance.
func CreateSimulation() *Simulation {
	var newSimulation Simulation = Simulation{}

	return &newSimulation
}

// AddIdle adds a single idle to the Simulation.
// Idle means two records that have been paired, but no action will be carried out because the relevant data is the same.
func (s *Simulation) AddIdle(pair PairSimulator) (bool, error) {
	var asStr string = pair.SimIdleString()
	fmt.Print(asStr)
	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert(pair PairSimulator) (bool, error) {
	var asStr string = pair.SimInsertString()
	fmt.Print(asStr)
	return false, nil
}

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(pair PairSimulator) (bool, error) {
	var asStr string = pair.SimUpdateString()
	fmt.Print(asStr)
	return false, nil
}
