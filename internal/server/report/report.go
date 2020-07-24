package report

import (
	"encoding/json"
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/synch"
)

type link struct {
	Cmd     string   `json:"cmd"`
	Idle    []string `json:"idle"`
	Inserts []string `json:"inserts"`
	Updates []string `json:"updates"`
}

// report is basically a report about what will happen after an actual synchronization is launched.
type report struct {
	msg   string
	synch *synch.Synch
	links map[string]*link
}

// SynchReportError is a custom report error.
type SynchReportError struct {
	SynchName string `json:"synchName"`
	ErrMsg    string `json:"errMsg"`
}

func (e *SynchReportError) Error() string {
	return fmt.Sprintf("['%s' synch report] %s", e.SynchName, e.ErrMsg)
}

// CreateReport creates a report instance.
func CreateReport(s *synch.Synch) *report {
	var newReport report = report{
		synch: s,
		links: make(map[string]*link),
	}

	return &newReport
}

// MarshalJSON implements the Marshaler interface for custom JSON creation.
func (r *report) MarshalJSON() ([]byte, error) {
	linkMap := make(map[int]*link)
	linkCounter := 1
	for _, link := range r.links {
		linkMap[linkCounter] = link
		linkCounter++
	}

	customStruct := struct {
		Msg       string           `json:"msg"`
		SynchInfo *cfg.SynchConfig `json:"synchInfo"`
		Links     map[int]*link    `json:"links"`
	}{
		Msg:       r.msg,
		SynchInfo: r.synch.GetConfig(),
		Links:     linkMap,
	}

	marshalled, err := json.Marshal(&customStruct)
	if err != nil {
		return nil, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
	}

	return marshalled, nil
}

// ToJSON turns the report into a JSON object.
func (r *report) ToJSON() ([]byte, error) {
	marshalled, err := json.Marshal(r)
	if err != nil {
		return nil, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
	}

	return marshalled, nil
}
