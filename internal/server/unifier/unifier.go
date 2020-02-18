package unifier

// Synchronizer is implemented by structs that do the actual synchronization actions.
type Synchronizer interface {
	Synchronize() (bool, error)
}

// Simulator enables creating strings which represent data flows
// needed for the simulation directly in the synch package (avoiding a package cycle).
type Simulator interface {
	AddIdle(pair Synchronizer) (bool, error)
	AddInsert(pair Synchronizer) (bool, error)
	AddUpdate(pair Synchronizer) (bool, error)
	Init()
	ToJSON() ([]byte, error)
}
