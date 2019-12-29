package synch

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var nullableFields = []string{"externalIds", "active"}

var synchConnectionTypes = [2]string{"external_id_columns", "persistence"}
var createNewRows = [3]string{"never", "initially", "always"}
var updateOldRows = [3]string{"never", "initially", "always"}

type synchData struct {
	Name    string   `yaml:"name"`
	Vectors []vector `yaml:"vectors"`
}

func (s *synchData) Validate() {
	validationUtil.YAMLStruct(*s, nullableFields)

	for _, vector := range s.Vectors {
		validationUtil.YAMLStruct(vector, nullableFields)

		validationUtil.YAMLStruct(vector.Source, nullableFields)
		validationUtil.YAMLStruct(vector.Source.Options, nullableFields)
		validationUtil.YAMLStruct(vector.Target, nullableFields)
		validationUtil.YAMLStruct(vector.Target.Options, nullableFields)
		validationUtil.YAMLStruct(vector.Settings, nullableFields)
		validationUtil.YAMLStruct(vector.Settings.ExternalIds, nullableFields)
	}
}
