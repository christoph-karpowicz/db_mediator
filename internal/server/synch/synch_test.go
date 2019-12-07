package synch

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestJSON(t *testing.T) {
	// os.Chdir("../../..")
	// synchs := Synchs{SynchMap: make(map[string]*Synch)}
	// synchs.ImportJSON()
}

func TestDir(t *testing.T) {
	files, err := ioutil.ReadDir("../../../config/synch-configs")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
