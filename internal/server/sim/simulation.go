package sim

type Simulation struct {
	updates []string
	inserts []string
}

func CreateSimulation() *Simulation {
	var newSimulation Simulation = Simulation{}

	return &newSimulation
}

func (s *Simulation) AddUpdate() (bool, error) {
	return false, nil
}

func (s *Simulation) AddInsert() (bool, error) {
	return false, nil
}
