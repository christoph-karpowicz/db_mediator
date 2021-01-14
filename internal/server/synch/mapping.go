package synch

import "github.com/christoph-karpowicz/unifier/internal/server/cfg"

// Mapping represents a single mapping in the config file like:
// example_node1.example_column1 TO example_node2.example_column2
type Mapping struct {
	synch        *Synch
	source       *node
	target       *node
	sourceColumn string
	targetColumn string
	raw          map[string]string
}

func createMapping(synch *Synch, mapping map[string]string) *Mapping {

	_, sourceNodeFound := mapping[cfg.PSUBEXP_SOURCE_NODE]
	if !sourceNodeFound {
		panic("[create mapping] ERROR: source node not found.")
	}
	_, targetNodeFound := mapping[cfg.PSUBEXP_TARGET_NODE]
	if !targetNodeFound {
		panic("[create mapping] ERROR: target node not found.")
	}

	newMapping := Mapping{
		synch:        synch,
		source:       synch.dbStore.nodes[mapping[cfg.PSUBEXP_SOURCE_NODE]],
		target:       synch.dbStore.nodes[mapping[cfg.PSUBEXP_TARGET_NODE]],
		sourceColumn: mapping[cfg.PSUBEXP_SOURCE_COLUMN],
		targetColumn: mapping[cfg.PSUBEXP_TARGET_COLUMN],
		raw:          mapping,
	}

	return &newMapping
}
