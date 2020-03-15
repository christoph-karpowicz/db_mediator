package unifier

// Synchronizer is implemented by structs that do the actual synchronization actions.
type Synchronizer interface {
	Synchronize() (bool, error)
}

// Reporter enables creating strings which represent data flows
// needed for the simulation directly in the synch package (avoiding a package cycle).
type Reporter interface {
	AddAction(p Synchronizer, actionType string) (bool, error)
	Finalize() ([]byte, error)
	Init()
}
