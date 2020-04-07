package synch

// Instruction represents an instruction from the config file.
type Instruction struct {
	synch    *Synch
	mappings []*Mapping
	links    []*Link
}

func createInstruction(synch *Synch) *Instruction {
	mpngs := make([]*Mapping, 0)
	links := make([]*Link, 0)

	newInstruction := Instruction{
		mappings: mpngs,
		links:    links,
		synch:    synch,
	}

	return &newInstruction
}
