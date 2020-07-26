package synch

import (
	"encoding/json"

	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

type History struct {
	rep     unifier.Reporter
	actions []*action
}

func (h *History) Init() {
	h.rep.Init()
}

// SetReporter inects a Reporter to the struct.
func (h *History) SetReporter(rep unifier.Reporter) {
	h.rep = rep
}

func (h *History) addAction(act *action) {
	h.actions = append(h.actions, act)
}

func (h History) GenerateReport() ([]byte, error) {
	err := h.addActionsToReport()
	if err != nil {
		panic(err)
	}

	report, err := h.rep.Finalize()
	if err != nil {
		panic(err)
	}
	return report, nil
}

func (h *History) addActionsToReport() error {
	for _, act := range h.actions {
		actionJSON, err := json.Marshal(&act)
		if err != nil {
			return err
			// return false, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
		}
		h.rep.AddAction(act.linkId, actionJSON, act.ActType)
	}
	return nil
}
