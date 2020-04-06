package synch

// Mapping represents a single mapping in the config file like:
// example_node1.example_column1 TO example_node2.example_column2
type Mapping struct {
	in           *Instruction
	source       *node
	target       *node
	sourceColumn string
	targetColumn string
}

func createMapping(in *Instruction, mapping map[string]string) *Mapping {

	_, sourceNodeFound := mapping["sourceNode"]
	if !sourceNodeFound {
		panic("[create mapping] ERROR: source node not found.")
	}
	_, targetNodeFound := mapping["targetNode"]
	if !targetNodeFound {
		panic("[create mapping] ERROR: target node not found.")
	}

	newMapping := Mapping{
		in:           in,
		source:       in.synch.nodes[mapping["sourceNode"]],
		target:       in.synch.nodes[mapping["targetNode"]],
		sourceColumn: mapping["sourceColumn"],
		targetColumn: mapping["targetColumn"],
	}

	return &newMapping
}
