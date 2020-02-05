package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{}

var synchConnectionTypes = [2]string{"external_id_columns", "persistence"}
var createNewRows = [3]string{"never", "initially", "always"}
var updateOldRows = [3]string{"never", "initially", "always"}

type synchData struct {
	Name     string     `yaml:"name"`
	Nodes    []nodeData `yaml:"nodes"`
	Mappings []string   `yaml:"mappings"`
}

func (s *synchData) Validate() {
	validationUtil.YAMLStruct(*s, nullableFields)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, nullableFields)
	}
}
