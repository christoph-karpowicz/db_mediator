package synch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Synchs struct {
	SynchMap map[string]*Synch
}

func (s *Synchs) ImportJSON() {
	synchFile, err := os.Open("config/synch.json")
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsing synch.json...")
	defer synchFile.Close()

	byteArray, err := ioutil.ReadAll(synchFile)
	if err != nil {
		panic(err)
	}

	var synchs map[string]json.RawMessage
	var synchsArray []Data

	json.Unmarshal(byteArray, &synchs)
	json.Unmarshal(synchs["synch"], &synchsArray)

	fmt.Println("----------------")
	fmt.Println("Synchs:")
	for i := 0; i < len(synchsArray); i++ {

		// Convert string arrays of vestors into Vector arrays.
		for j := 0; j < len(synchsArray[i].Tables); j++ {
			for k := 0; k < len(synchsArray[i].Tables[j].Vectors); k++ {
				pair := synchsArray[i].Tables[j].Vectors[k].(string)
				synchsArray[i].Tables[j].Vectors[k] = Vector{}
				v := synchsArray[i].Tables[j].Vectors[k].(Vector)
				vPtr := &v
				vPtr.Parse(&pair)
			}
		}

		s.SynchMap[synchsArray[i].Name] = &Synch{synch: &synchsArray[i]}
		fmt.Println(synchsArray[i])

	}
	fmt.Println("----------------")
}
