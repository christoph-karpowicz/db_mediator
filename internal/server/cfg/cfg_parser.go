package cfg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ImportYAMLDir invokes the import function on all .yaml files from
// a directory.
func ImportYAMLDir(dirPath string) []*SynchConfig {
	var synchCfgArray []*SynchConfig = make([]*SynchConfig, 0)

	configFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------------")
	fmt.Println("Synchs:")
	for _, configFile := range configFiles {
		var synchCfg SynchConfig = SynchConfig{}
		ImportYAMLFile(&synchCfg, dirPath+"/"+configFile.Name())

		// Don't load inactive synchs.
		// if !synchData.Settings.Active {
		// 	continue
		// }

		synchCfgArray = append(synchCfgArray, &synchCfg)
		fmt.Println("Config file name: " + configFile.Name())
		// fmt.Println(synchData)
	}
	fmt.Println("----------------")

	return synchCfgArray
}

// ImportYAMLFile imports a configuration file into a Config struct.
func ImportYAMLFile(cfg Config, filePath string) {
	fp, _ := filepath.Abs(filePath)

	configFile, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsing %s...\n", fp)
	defer configFile.Close()

	byteArray, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	var marshalErr error

	switch cfg.(type) {
	case *DbConfigArray:
		marshalErr = yaml.Unmarshal(byteArray, cfg.(*DbConfigArray))
	case *SynchConfig:
		marshalErr = yaml.Unmarshal(byteArray, cfg.(*SynchConfig))
	}

	if marshalErr != nil {
		log.Fatalf("error: %v", marshalErr)
	}
}
