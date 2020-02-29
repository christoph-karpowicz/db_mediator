package synch

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Synchs is a collection of all synchronizations.
type Synchs struct {
	SynchMap map[string]*Synch
}

// ImportYAMLDir invokes the import function on all .yaml files from
// a directory.
func (s *Synchs) ImportYAMLDir() {
	configFiles, err := ioutil.ReadDir("./config/synchs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------------")
	fmt.Println("Synchs:")
	for _, configFile := range configFiles {
		synchData := s.ImportYAMLFile(configFile.Name())

		// Don't load inactive synchs.
		// if !synchData.Settings.Active {
		// 	continue
		// }

		s.SynchMap[synchData.Name] = &Synch{Cfg: &synchData, initial: true}
		fmt.Println("Config file name: " + configFile.Name())
		// fmt.Println(synchData)
	}
	fmt.Println("----------------")

}

// ImportYAMLFile imports a configuration file into a Config struct.
func (s *Synchs) ImportYAMLFile(fileName string) Config {
	synchFilePath, _ := filepath.Abs("./config/synchs/" + fileName)

	synchFile, err := os.Open(synchFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsing %s...\n", synchFilePath)
	defer synchFile.Close()

	byteArray, err := ioutil.ReadAll(synchFile)
	if err != nil {
		panic(err)
	}

	var synch Config = Config{}

	marshalErr := yaml.Unmarshal(byteArray, &synch)
	if marshalErr != nil {
		log.Fatalf("error: %v", marshalErr)
	}

	return synch
}

// ValidateYAML validates data imported from a config file.
func (s *Synchs) ValidateYAML() {
	fmt.Println("Synch YAML file validation...")
	for _, synch := range s.SynchMap {
		(*synch).GetConfig().Validate()
	}
	fmt.Println("... passed.")
}

// CreateSynchs constructor function for the Synchs struct.
func CreateSynchs() *Synchs {
	synchs := &Synchs{SynchMap: make(map[string]*Synch)}
	return synchs
}
