package synch

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestYAML(t *testing.T) {
	os.Chdir("../../..")
	synchs := Synchs{SynchMap: make(map[string]*Synch)}
	synchs.ImportYAMLDir()
	synchs.ValidateYAML()
	// synchs.SynchMap["films"].Init()
	// synchs.SynchMap["films"].Synchronize()
}

func TestDir(t *testing.T) {
	files, err := ioutil.ReadDir("./config/synchs")
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		log.Fatal("no config files")
	}
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// }
}
