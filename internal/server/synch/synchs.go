package synch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Synchs struct {
	SynchMap map[string]*Synch
}

func (s *Synchs) ImportJSONDir() {
	configFiles, err := ioutil.ReadDir("./config/synch-configs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------------")
	fmt.Println("Synchs:")
	for _, configFile := range configFiles {
		synchData := s.ImportJSONFile(configFile.Name())

		// Don't load inactive synchs.
		if !synchData.Active {
			continue
		}

		s.SynchMap[synchData.Name] = &Synch{synch: &synchData, initial: true}
		fmt.Println("Config file name: " + configFile.Name())
		fmt.Println(synchData)
	}
	fmt.Println("----------------")

}

func (s *Synchs) ImportJSONFile(fileName string) SynchData {
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

	var synch SynchData

	json.Unmarshal(byteArray, &synch)

	return synch
}

func (s *Synchs) ValidateJSON() {
	fmt.Println("Synch JSON file validation...")
	for _, synch := range s.SynchMap {
		fmt.Println((*synch).GetData())
		(*synch).GetData().Validate()
	}
	fmt.Println("...passed.")
}
