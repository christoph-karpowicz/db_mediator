package report

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/synch"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// Reporter acts as an interface for other packages
// to modify the synch report.
type Reporter struct {
	synch *synch.Synch
	rep   *report
}

// CreateReporter creates a Reporter instance.
func CreateReporter(s *synch.Synch) unifier.Reporter {
	var newReporter Reporter = Reporter{
		s,
		CreateReport(s),
	}

	return &newReporter
}

// Init fills the necessary fields after the Synch instance finished its Init execution.
func (r *Reporter) Init() {
	for _, lnk := range r.synch.Links {
		_, lnkExists := r.rep.links[lnk.GetID()]
		if !lnkExists {
			r.rep.links[lnk.GetID()] = &link{Cmd: lnk.Cmd}
		}
	}
}

// AddAction adds a single action to the report.
// Action types:
// 	1.	idle - means two records that have been paired, but no action will be carried out because the relevant data is the same.
// 	2.	insert
// 	3. 	update
func (r *Reporter) AddAction(linkID string, actionJSON []byte, actionType string) (bool, error) {
	switch actionType {
	case cfg.OPERATION_IDLE:
		r.rep.links[linkID].Idle = append(r.rep.links[linkID].Idle, string(actionJSON))
	case cfg.OPERATION_INSERT:
		r.rep.links[linkID].Inserts = append(r.rep.links[linkID].Inserts, string(actionJSON))
	case cfg.OPERATION_UPDATE:
		r.rep.links[linkID].Updates = append(r.rep.links[linkID].Updates, string(actionJSON))
	}

	return true, nil
}

func (r *Reporter) SetReportMessage(msg string) {
	r.rep.msg = msg
}

// Finalize wraps up the report creation process.
func (r *Reporter) Finalize() ([]byte, error) {
	toJSON, err := r.rep.ToJSON()
	if err != nil {
		return nil, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
	}

	return toJSON, nil
}
