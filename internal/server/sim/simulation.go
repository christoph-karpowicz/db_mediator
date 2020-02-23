package sim

import (
	"encoding/json"
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
	Idle    []string `json:"idle"`
	Inserts []string `json:"inserts"`
	Updates []string `json:"updates"`
}

type mappingSimulation struct {
	mapping  *synch.Mapping
	LinkSims map[int]*linkSimulation `json:"links"`
}

// Simulation is basically a report about what will happen after an actual synchronization is launched.
type Simulation struct {
	synch        *synch.Synch
	mappingsSims map[*synch.Mapping]*mappingSimulation
}

// CreateSimulation creates a Simulation instance.
func CreateSimulation(s *synch.Synch) unifier.Simulator {
	var newSimulation Simulation = Simulation{
		synch:        s,
		mappingsSims: make(map[*synch.Mapping]*mappingSimulation),
	}

	return &newSimulation
}

// AddIdle adds a single idle to the Simulation.
// Idle means two records that have been paired, but no action will be carried out because the relevant data is the same.
func (s *Simulation) AddIdle(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimIdleString()
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Idle = append(s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Idle, str)
	fmt.Print(str)

	return false, nil
}

// AddInsert adds a single insert to the Simulation.
func (s *Simulation) AddInsert(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimInsertString()
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Inserts = append(s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Inserts, str)
	fmt.Print(str)

	return false, nil
}

// AddUpdate adds a single update to the Simulation.
func (s *Simulation) AddUpdate(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.SimUpdateString()
	var lnkIdx int = pair.Mapping.Sim.LinkIndex

	s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Updates = append(s.mappingsSims[pair.Mapping].LinkSims[lnkIdx].Updates, str)
	fmt.Print(str)

	return false, nil
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (s *Simulation) Init() {
	for _, mapping := range s.synch.Mappings {
		_, mpngSimExists := s.mappingsSims[mapping]
		if !mpngSimExists {
			s.mappingsSims[mapping] = &mappingSimulation{mapping: mapping}
			s.mappingsSims[mapping].LinkSims = make(map[int]*linkSimulation)
		}

		s.mappingsSims[mapping].LinkSims[mapping.Sim.LinkIndex] = &linkSimulation{}
	}
}

// MarshalJSON implements the Marshaler interface for custom JSON creation.
func (s *Simulation) MarshalJSON() ([]byte, error) {
	mappingsMap := make(map[int]*mappingSimulation)
	for _, mappingSim := range s.mappingsSims {
		_, mpngSimExists := mappingsMap[mappingSim.mapping.Sim.MappingIndex]
		if !mpngSimExists {
			mappingsMap[mappingSim.mapping.Sim.MappingIndex] = mappingSim
		} else {
			for k, v := range mappingSim.LinkSims {
				mappingsMap[mappingSim.mapping.Sim.MappingIndex].LinkSims[k] = v
			}
		}
	}

	customStruct := struct {
		SynchInfo    *synch.Config              `json:"synchInfo"`
		MappingsSims map[int]*mappingSimulation `json:"mappings"`
	}{
		SynchInfo:    s.synch.Cfg,
		MappingsSims: mappingsMap,
	}

	return json.Marshal(&customStruct)
}

// ToJSON turns the simulation into a JSON object.
func (s *Simulation) ToJSON() ([]byte, error) {
	fmt.Println(s)
	return json.Marshal(s)
}
