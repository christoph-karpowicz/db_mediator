package synch

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

type History struct {
	synch   *Synch
	rep     unifier.Reporter
	actions []*operation
}

func (h *History) Init(synch *Synch) {
	h.synch = synch
	h.rep.Init()
}

// SetReporter inects a Reporter to the struct.
func (h *History) SetReporter(rep unifier.Reporter) {
	h.rep = rep
}

func (h *History) addAction(act *operation) {
	h.actions = append(h.actions, act)
}

func (h History) GenerateReport() ([]byte, error) {
	err := h.addActionsToReport()
	if err != nil {
		panic(err)
	}

	h.setReportMessage()
	report, err := h.rep.Finalize()
	if err != nil {
		panic(err)
	}
	return report, nil
}

func (h *History) addActionsToReport() error {
	// for _, act := range h.actions {
	// 	actionJSON, err := json.Marshal(&act)
	// 	if err != nil {
	// 		return err
	// 		// return false, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
	// 	}
	// 	h.rep.AddAction(act.linkId, actionJSON, act.ActType)
	// }
	return nil
}

func (h *History) setReportMessage() {
	fmt.Println(h.synch.stype)
	if h.synch.IsSimulation() {
		h.rep.SetReportMessage("'" + h.synch.GetConfig().Name + "' simulation was successful. " +
			"The report contains changes that would be made if you ran an actual synchronization.")
	} else if h.synch.stype == ONGOING {
		h.rep.SetReportMessage("Synch '" + h.synch.GetConfig().Name + "' stopped. " +
			"The report contains changes that have been made to the relevant nodes.")
	} else {
		h.rep.SetReportMessage("'" + h.synch.GetConfig().Name + "' synchronization was successful. " +
			"The report contains changes that have been made to the relevant nodes.")
	}
}
