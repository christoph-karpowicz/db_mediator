package synch

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var synchs Synchs

func TestYAML(t *testing.T) {
	os.Chdir("../../..")
	synchs = CreateSynchs()
	synchs.Init()
}

func TestDir(t *testing.T) {
	files, err := ioutil.ReadDir("./config/synch")
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
