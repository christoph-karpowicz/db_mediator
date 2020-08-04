package cfg

import (
	"fmt"

	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

var watcherNullableFields = []string{"name", "key"}

// WatcherConfig holds raw data from the YAML config file.
type WatcherConfig struct {
	Name  string       `yaml:"name"`
	Nodes []NodeConfig `yaml:"nodes"`
}

// Validate data from the YAML file.
func (s *WatcherConfig) Validate() {
	validationUtil.YAMLStruct(*s, nil)

	for _, node := range s.Nodes {
		validationUtil.YAMLStruct(node, watcherNullableFields)
	}
}

// GetWatcherConfigs loads configs from the synchs directory.
func GetWatcherConfigs() []Config {
	fmt.Println("Watchers:")
	return ImportYAMLDir(WATCHER_DIR)
}
