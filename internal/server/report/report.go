package report

import (
	"encoding/json"
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/synch"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

type linkRep struct {
	Cmd     string   `json:"cmd"`
	Idle    [][]byte `json:"idle"`
	Inserts [][]byte `json:"inserts"`
	Updates [][]byte `json:"updates"`
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

// ReportError is a custom report error.
type ReportError struct {
	SynchName string `json:"synchName"`
	ErrMsg    string `json:"errMsg"`
}

func (e *ReportError) Error() string {
	return fmt.Sprintf("['%s' synch report] %s", e.SynchName, e.ErrMsg)
}

// CreateReport creates a Report instance.
func CreateReport(s *synch.Synch) unifier.Reporter {
	var newReport Report = Report{
		synch:        s,
		mappingsReps: make(map[*synch.Mapping]*mappingRep),
	}

	return &newReport
}

// AddAction adds a single action to the Report.
// Action types:
// 	1.	idle - means two records that have been paired, but no action will be carried out because the relevant data is the same.
// 	2.	insert
// 	3. 	update
func (r *Report) AddAction(p unifier.Synchronizer, actionType string) (bool, error) {
	var pair synch.Pair = p.(synch.Pair)
	var lnkIdx int = pair.Mapping.Rep.LinkIndex
	actionJSON, err := pair.ReportJSON(actionType)
	if err != nil {
		return false, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	switch actionType {
	case "idle":
		r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle, actionJSON)
	case "insert":
		r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Inserts = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle, actionJSON)
	case "update":
		r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Updates = append(r.mappingsReps[pair.Mapping].LinkReps[lnkIdx].Idle, actionJSON)
	}
	// fmt.Print(actionJSON)

	return true, nil
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

	toJSON, err := r.ToJSON()
	if err != nil {
		return nil, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	return toJSON, nil
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

	marshal, err := json.Marshal(&customStruct)
	if err != nil {
		return nil, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	return marshal, nil
}

// ToJSON turns the report into a JSON object.
func (r *Report) ToJSON() ([]byte, error) {
	// fmt.Println(s)
	marshal, err := json.Marshal(r)
	if err != nil {
		return nil, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	return marshal, nil
}
