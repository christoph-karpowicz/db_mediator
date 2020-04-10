package report

import (
	"encoding/json"
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/synch"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

type linkRep struct {
	Cmd     string   `json:"cmd"`
	Idle    []string `json:"idle"`
	Inserts []string `json:"inserts"`
	Updates []string `json:"updates"`
}

type instructionRep struct {
	inPtr      *synch.Instruction
	LnkReports map[int]*linkRep `json:"links"`
}

// Report is basically a report about what will happen after an actual synchronization is launched.
type Report struct {
	msg       string
	synch     *synch.Synch
	inReports map[*synch.Instruction]*instructionRep
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
		synch:     s,
		inReports: make(map[*synch.Instruction]*instructionRep),
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
	var lnkIdx int = pair.Link.Rep.LinkIndex
	actionJSON, err := pair.ReportJSON(actionType)
	if err != nil {
		return false, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	switch actionType {
	case "idle":
		r.inReports[pair.Link.In].LnkReports[lnkIdx].Idle = append(r.inReports[pair.Link.In].LnkReports[lnkIdx].Idle, string(actionJSON))
	case "insert":
		r.inReports[pair.Link.In].LnkReports[lnkIdx].Inserts = append(r.inReports[pair.Link.In].LnkReports[lnkIdx].Inserts, string(actionJSON))
	case "update":
		r.inReports[pair.Link.In].LnkReports[lnkIdx].Updates = append(r.inReports[pair.Link.In].LnkReports[lnkIdx].Updates, string(actionJSON))
	}
	// fmt.Print(actionJSON)

	return true, nil
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (r *Report) Init() {
	for _, in := range r.synch.Instructions {
		_, inRepExists := r.inReports[in]
		if !inRepExists {
			r.inReports[in] = &instructionRep{inPtr: in}
			r.inReports[in].LnkReports = make(map[int]*linkRep)
		}

		for _, lnk := range in.Links {
			r.inReports[in].LnkReports[lnk.Rep.LinkIndex] = &linkRep{Cmd: lnk.Rep.Link["raw"]}
		}

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
// func (r *Report) MarshalJSON() ([]byte, error) {
// 	instructionsMap := make(map[int]*instructionRep)
// 	for _, instructionRep := range r.inReports {
// 		_, inRepExists := instructionsMap[instructionRep.inPtr.Rep.LinkIndex]
// 		if !inRepExists {
// 			instructionsMap[instructionRep.inPtr.Rep.LinkIndex] = instructionRep
// 		} else {
// 			for k, v := range instructionRep.LnkReports {
// 				instructionsMap[instructionRep.inPtr.Rep.LinkIndex].LnkReports[k] = v
// 			}
// 		}
// 	}

// 	customStruct := struct {
// 		Msg       string                  `json:"msg"`
// 		SynchInfo *synch.Config           `json:"synchInfo"`
// 		InReps    map[int]*instructionRep `json:"instructions"`
// 	}{
// 		Msg:       r.msg,
// 		SynchInfo: r.synch.Cfg,
// 		InReps:    instructionsMap,
// 	}

// 	marshalled, err := json.Marshal(&customStruct)
// 	if err != nil {
// 		return nil, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
// 	}

// 	return marshalled, nil
// }

// ToJSON turns the report into a JSON object.
func (r *Report) ToJSON() ([]byte, error) {
	// fmt.Println(s)
	marshalled, err := json.Marshal(r)
	if err != nil {
		return nil, &ReportError{SynchName: r.synch.Cfg.Name, ErrMsg: err.Error()}
	}

	return marshalled, nil
}
