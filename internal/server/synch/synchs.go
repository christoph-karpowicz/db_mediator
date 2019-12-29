/*
Package synch handles all data sychronization.
*/
package synch

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Synchs struct {
	SynchMap map[string]*synch
}

func (s *Synchs) ImportYAMLDir() {
	configFiles, err := ioutil.ReadDir("./config/synch-configs")
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

		s.SynchMap[synchData.Name] = &synch{synch: &synchData, initial: true}
		fmt.Println("Config file name: " + configFile.Name())
		fmt.Println(synchData)
	}
	fmt.Println("----------------")

}

func (s *Synchs) ImportYAMLFile(fileName string) synchData {
	synchFilePath, _ := filepath.Abs("./config/synch-configs/" + fileName)

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

	var synch synchData = synchData{}

	marshalErr := yaml.Unmarshal(byteArray, &synch)
	if marshalErr != nil {
		log.Fatalf("error: %v", marshalErr)
	}

	return synch
}

func (s *Synchs) ValidateYAML() {
	fmt.Println("Synch YAML file validation...")
	for _, synch := range s.SynchMap {
		fmt.Println((*synch).GetData())
		(*synch).GetData().Validate()
	}
	fmt.Println("...passed.")
}

func CreateSynchs() *Synchs {
	synchs := &Synchs{SynchMap: make(map[string]*synch)}
	return synchs
}
