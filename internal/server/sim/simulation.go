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

type linkSimulation struct {
	idle    []string
	inserts []string
	updates []string
}

type mappingSimulation struct {
	mapping  *synch.Mapping
	linkSims map[int]*linkSimulation
}

// Simulation is basically a report about what will happen after an actual synchronization is launched.
type Simulation struct {
	synch        *synch.Synch
	mappingsSims map[int]*mappingSimulation
}

// CreateSimulation creates a Simulation instance.
func CreateSimulation(s *synch.Synch) unifier.Simulator {
	var newSimulation Simulation = Simulation{
		synch:        s,
		mappingsSims: make(map[int]*mappingSimulation),
	}

	return &newSimulation
}

// AddIdle adds a single idle to the Simulation.
// Idle means two records that have been paired, but no action will be carried out because the relevant data is the same.
func (s *Simulation) AddIdle(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimIdleString()
	var mpngIdx int = pair.Mapping.Sim.MappingIndex
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[mpngIdx].linkSims[lnkIdx].idle = append(s.mappingsSims[mpngIdx].linkSims[lnkIdx].idle, str)
	fmt.Print(str)

	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimInsertString()
	var mpngIdx int = pair.Mapping.Sim.MappingIndex
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[mpngIdx].linkSims[lnkIdx].inserts = append(s.mappingsSims[mpngIdx].linkSims[lnkIdx].inserts, str)
	fmt.Print(str)

	return false, nil
}

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimUpdateString()
	var mpngIdx int = pair.Mapping.Sim.MappingIndex
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[mpngIdx].linkSims[lnkIdx].updates = append(s.mappingsSims[mpngIdx].linkSims[lnkIdx].updates, str)
	fmt.Print(str)

	return false, nil
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (s *Simulation) Init() {
	for _, mapping := range s.synch.Mappings {
		_, mpngSimExists := s.mappingsSims[mapping.Sim.MappingIndex]
		if !mpngSimExists {
			s.mappingsSims[mapping.Sim.MappingIndex] = &mappingSimulation{mapping: mapping}
			s.mappingsSims[mapping.Sim.MappingIndex].linkSims = make(map[int]*linkSimulation)
		}

		s.mappingsSims[mapping.Sim.MappingIndex].linkSims[mapping.Sim.LinkIndex] = &linkSimulation{}
	}
}
