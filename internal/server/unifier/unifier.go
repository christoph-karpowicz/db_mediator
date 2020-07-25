package unifier

// Reporter enables creating strings which represent data flows
// needed for the simulation directly in the synch package (avoiding a package cycle).
type Reporter interface {
	AddAction(linkID string, actionJSON []byte, actionType string) (bool, error)
	Finalize() ([]byte, error)
	Init()
}

// Pairable is used in generating reports to get pair
// data.
type Pairable interface {
	ReportJSON(actionType string) ([]byte, error)
	GetLinkID() string
}
