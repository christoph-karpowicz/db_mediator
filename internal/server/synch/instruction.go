package synch

// Instruction represents an instruction from the config file.
type Instruction struct {
	synch    *Synch
	mappings []*Mapping
	Links    []*Link
}

func createInstruction(synch *Synch) *Instruction {
	mpngs := make([]*Mapping, 0)
	Links := make([]*Link, 0)

	newInstruction := Instruction{
		mappings: mpngs,
		Links:    Links,
		synch:    synch,
	}

	return &newInstruction
}
