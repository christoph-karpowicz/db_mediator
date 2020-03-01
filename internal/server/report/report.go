package report

import (
	"encoding/json"

	"github.com/christoph-karpowicz/unifier/internal/server/synch"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// PairReporter is implemented by the pair struct from the synch package
// and passed to Add* methods.
type PairReporter interface {
	RepIdleString() string
	RepInsertString() string
	RepUpdateString() string
}

type linkRep struct {
	Cmd     string   `json:"cmd"`
	Idle    []string `json:"idle"`
	Inserts []string `json:"inserts"`
	Updates []string `json:"updates"`
}

type mappingRep struct {
	mpngPtr  *synch.Mapping
	LinkReps map[int]*linkRep `json:"links"`
}

// Report is basically a report about what will happen after an actual synchronization is launched.
type Report struct {
	synch        *synch.Synch
	mappingsReps map[*synch.Mapping]*mappingRep
}

// CreateReport creates a Report instance.
func CreateReport(s *synch.Synch) unifier.Reporter {
	var newReport Report = Report{
		synch:        s,
		mappingsReps: make(map[*synch.Mapping]*mappingRep),
	}

	return &newReport
}

// AddIdle adds a single idle to the Report.
// Idle means two records that have been paired, but no action will be carried out because the relevant data is the same.
func (s *Report) AddIdle(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepIdleString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle = append(s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle, str)
	// fmt.Print(str)

	return false, nil
}

// AddInsert adds a single insert to the Report.
func (s *Report) AddInsert(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepInsertString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Inserts = append(s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Inserts, str)
	// fmt.Print(str)

	return false, nil
}

// AddUpdate adds a single update to the Report.
func (s *Report) AddUpdate(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepUpdateString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Updates = append(s.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Updates, str)
	// fmt.Print(str)

	return false, nil
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (s *Report) Init() {
	for _, mpng := range s.synch.Mappings {
		_, mpngRepExists := s.mappingsReps[mpng]
		if !mpngRepExists {
			s.mappingsReps[mpng] = &mappingRep{mpngPtr: mpng}
			s.mappingsReps[mpng].LinkReps = make(map[int]*linkRep)
		}

		s.mappingsReps[mpng].LinkReps[mpng.Rep.LinkIndex] = &linkRep{Cmd: mpng.Rep.Link["raw"]}
	}
}

// MarshalJSON implements the Marshaler interface for custom JSON creation.
func (s *Report) MarshalJSON() ([]byte, error) {
	mappingsMap := make(map[int]*mappingRep)
	for _, mappingRep := range s.mappingsReps {
		_, mpngRepExists := mappingsMap[mappingRep.mpngPtr.Rep.MappingIndex]
		if !mpngRepExists {
			mappingsMap[mappingRep.mpngPtr.Rep.MappingIndex] = mappingRep
		} else {
			for k, v := range mappingRep.LinkReps {
				mappingsMap[mappingRep.mpngPtr.Rep.MappingIndex].LinkReps[k] = v
			}
		}
	}

	customStruct := struct {
		SynchInfo   *synch.Config       `json:"synchInfo"`
		MappingReps map[int]*mappingRep `json:"mappings"`
	}{
		SynchInfo:   s.synch.Cfg,
		MappingReps: mappingsMap,
	}

	return json.Marshal(&customStruct)
}

// ToJSON turns the report into a JSON object.
func (s *Report) ToJSON() ([]byte, error) {
	// fmt.Println(s)
	return json.Marshal(s)
}
