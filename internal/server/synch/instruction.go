package synch

// InstructionReportData Instruction data for simulation purposes.
type InstructionReportData struct {
	InstructionIndex int
	LinkIndex        int
	Link             map[string]string
}

// Instruction represents an instruction from the config file.
type Instruction struct {
	synch    *Synch
	mappings []*Mapping
	links    []*Link
	Rep      *InstructionReportData
}

func createInstruction(synch *Synch) *Instruction {
	mpngs := make([]*Mapping, 0)
	links := make([]*Link, 0)

	newInstruction := Instruction{
		mappings: mpngs,
		links:    links,
		synch:    synch,
	}

	// newInstruction.Rep = &InstructionReportData{
	// 	InstructionIndex: indexes[0],
	// 	LinkIndex:        indexes[1],
	// 	Link:             link,
	// }

	return &newInstruction
}
