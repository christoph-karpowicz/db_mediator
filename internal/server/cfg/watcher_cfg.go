package cfg

import (
	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

// WatcherConfig holds raw data from the YAML config file.
type WatcherConfig struct {
	Name  string   `yaml:"name"`
	Nodes []string `yaml:"nodes"`
}

// Validate data from the YAML file.
func (s *WatcherConfig) Validate() {
	validationUtil.YAMLStruct(*s, synchNullableFields)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, synchNullableFields)
	}
}

// GetWatcherConfigs loads configs from the synchs directory.
func GetWatcherConfigs() []*WatcherConfig {
	return ImportYAMLDir("./config/watch")
}
