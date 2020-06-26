package cfg

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var synchNullableFields = []string{}

var synchConnectionTypes = [2]string{"external_id_columns", "persistence"}
var createNewRows = [3]string{"never", "initially", "always"}
var updateOldRows = [3]string{"never", "initially", "always"}

// Config holds raw data from the YAML config file.
type SynchConfig struct {
	Name    string       `yaml:"name"`
	Nodes   []NodeConfig `yaml:"nodes"`
	Map     []string     `yaml:"map"`
	Link    []string     `yaml:"link"`
	MatchBy MatchBy      `yaml:"match_by"`
	Do      []string     `yaml:"do"`
}

// Validate data from the YAML file.
func (s *SynchConfig) Validate() {
	validationUtil.YAMLStruct(*s, synchNullableFields)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, synchNullableFields)
	}
}

func GetSynchConfigs() []*SynchConfig {
	return ImportYAMLDir("./config/synchs")
}