package synch

import (
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	os.Chdir("../..")
	synchs := Synchs{SynchMap: make(map[string]*Synch)}
	synchs.ImportJSON()
}

func TestVectorParsing(t *testing.T) {
	var vector *Vector = &Vector{}
	example := "example1 <=> example2"
	vector.Parse(&example)
	parsedCorrectly := (vector.Column1 == "example1" && vector.Column2 == "example2" && vector.Dir == "<=>")
	if !parsedCorrectly {
		t.Error("Vector hasn't been parsed correctly.")
	}
}
