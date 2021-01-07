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
func ImportYAMLDir(dirPath string) []Config {
	var cfgs []Config = make([]Config, 0)

	configFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, configFile := range configFiles {
		var config Config
		switch dirPath {
		case SYNCH_DIR:
			config = &SynchConfig{}
			cfgs = append(cfgs, config)
		}

		fmt.Println("Config file name: " + configFile.Name())
		ImportYAMLFile(config, dirPath+"/"+configFile.Name())
	}

	return cfgs
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
