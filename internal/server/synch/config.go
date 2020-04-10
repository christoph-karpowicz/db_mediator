package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{}

var synchConnectionTypes = [2]string{"external_id_columns", "persistence"}
var createNewRows = [3]string{"never", "initially", "always"}
var updateOldRows = [3]string{"never", "initially", "always"}

// Config holds raw data from the YAML config file.
type Config struct {
	Name    string       `yaml:"name"`
	Nodes   []nodeConfig `yaml:"nodes"`
	Map     []string     `yaml:"map"`
	Synch   []string     `yaml:"synch"`
	MatchBy matchBy      `yaml:"match_by"`
	Do      []string     `yaml:"do"`
}

// Validate data from the YAML file.
func (s *Config) Validate() {
	validationUtil.YAMLStruct(*s, nullableFields)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, nullableFields)
	}
}
