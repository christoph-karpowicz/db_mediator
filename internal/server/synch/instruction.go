package synch

// Instruction represents an instruction from the config file.
type Instruction struct {
	synch        *Synch
	source       *node
	target       *node
	sourceColumn string
	targetColumn string
	sourceWhere  string
	targetWhere  string
}

func createInstruction(synch *Synch, in map[string]string) *Instruction {
	newInstruction := Instruction{
		synch:        synch,
		source:       synch.nodes[in["sourceNode"]],
		target:       synch.nodes[in["targetNode"]],
		sourceColumn: in["sourceColumn"],
		targetColumn: in["targetColumn"],
		sourceWhere:  in["sourceWhere"],
		targetWhere:  in["targetWhere"],
	}

	return &newInstruction
}
