package sim

import "fmt"

// Simulation is basically a report about what will happen after an actual synchronization is launched.
type Simulation struct {
	updates []string
	inserts []string
}

// CreateSimulation creates a Simulation instance.
func CreateSimulation() *Simulation {
	var newSimulation Simulation = Simulation{}

	return &newSimulation
}

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(source *RecordState, target *RecordState) (bool, error) {
	var asStr string = fmt.Sprintf("|%6v: %6v, %6v: %16v| => |%6v: %6v, %6s: %16v -> %16v|\n", source.KeyName, source.KeyValue, source.ColumnName, source.CurrentValue, target.KeyName, target.KeyValue, target.ColumnName, target.CurrentValue, target.NewValue)
	fmt.Println(asStr)
	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert() (bool, error) {
	return false, nil
}
