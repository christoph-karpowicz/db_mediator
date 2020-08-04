package cfg

import (
	"fmt"

	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var synchNullableFields = []string{}

const (
	DB_INSERT = "INSERT"
	DB_UPDATE = "UPDATE"
)

// SynchConfig holds raw data from the YAML config file.
type SynchConfig struct {
	Name  string       `yaml:"name"`
	Nodes []NodeConfig `yaml:"nodes"`
	Map   []string     `yaml:"map"`
	Link  []string     `yaml:"link"`
	Match Match        `yaml:"match"`
	Do    []string     `yaml:"do"`
}

// Validate data from the YAML file.
func (s *SynchConfig) Validate() {
	validationUtil.YAMLStruct(*s, synchNullableFields)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, synchNullableFields)
	}
}

// GetSynchConfigs loads configs from the synchs directory.
func GetSynchConfigs() []Config {
	fmt.Println("Synchs:")
	return ImportYAMLDir(SYNCH_DIR)
}
