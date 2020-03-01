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
	msg          string
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
func (r *Report) AddIdle(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepIdleString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle, str)
	// fmt.Print(str)

	return false, nil
}

// AddInsert adds a single insert to the Report.
func (r *Report) AddInsert(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepInsertString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Inserts = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Inserts, str)
	// fmt.Print(str)

	return false, nil
}

// AddUpdate adds a single update to the Report.
func (r *Report) AddUpdate(p unifier.Synchronizer) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var str string = pair.RepUpdateString()
	var lnkIdx int = pair.Mapping.Rep.LinkIndex

	r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Updates = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Updates, str)
	// fmt.Print(str)

	return false, nil
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (r *Report) Init() {
	for _, mpng := range r.synch.Mappings {
		_, mpngRepExists := r.mappingsReps[mpng]
		if !mpngRepExists {
			r.mappingsReps[mpng] = &mappingRep{mpngPtr: mpng}
			r.mappingsReps[mpng].LinkReps = make(map[int]*linkRep)
		}

		r.mappingsReps[mpng].LinkReps[mpng.Rep.LinkIndex] = &linkRep{Cmd: mpng.Rep.Link["raw"]}
	}
}

// Finalize wraps up the report creation process.
func (r *Report) Finalize() ([]byte, error) {
	if r.synch.Simulation {
		r.msg = "'" + r.synch.GetConfig().Name + "' simulation was successful. The report contains changes that would be made if you requested an actual synchronization."
	} else {
		r.msg = "'" + r.synch.GetConfig().Name + "' synchronization was successful. The report contains changes that have been made to the relevant nodes."
	}

	return r.ToJSON()
}

// MarshalJSON implements the Marshaler interface for custom JSON creation.
func (r *Report) MarshalJSON() ([]byte, error) {
	mappingsMap := make(map[int]*mappingRep)
	for _, mappingRep := range r.mappingsReps {
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
		Msg         string              `json:"msg"`
		SynchInfo   *synch.Config       `json:"synchInfo"`
		MappingReps map[int]*mappingRep `json:"mappings"`
	}{
		Msg:         r.msg,
		SynchInfo:   r.synch.Cfg,
		MappingReps: mappingsMap,
	}

	return json.Marshal(&customStruct)
}

// ToJSON turns the report into a JSON object.
func (r *Report) ToJSON() ([]byte, error) {
	// fmt.Println(s)
	return json.Marshal(r)
}
