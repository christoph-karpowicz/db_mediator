package sim

import "fmt"

// Simulator enables creating strings which represent data flows
// needed for the simulation directly in the synch package (avoiding a package cycle).
type Simulator interface {
	CreateSimulationString() string
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

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(pair Simulator) (bool, error) {
	var asStr string = pair.CreateSimulationString()
	fmt.Print(asStr)
	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert() (bool, error) {
	return false, nil
}
