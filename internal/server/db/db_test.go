package db

import (
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	os.Chdir("../../..")
	dbs := Databases{DBMap: make(map[string]*Database)}
	dbs.ImportJSON()
}
