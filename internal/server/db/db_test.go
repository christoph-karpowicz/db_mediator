package db

import (
	"os"
	"testing"
)

func TestYAML(t *testing.T) {
	os.Chdir("../../..")
	dbs := Databases{DBMap: make(map[string]*Database)}
	dbs.ImportYAML()
	dbs.ValidateYAML()
}
