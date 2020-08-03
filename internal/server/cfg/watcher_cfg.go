package cfg

import (
	"fmt"

	validationUtil "github.com/christoph-karpowicz/unifier/internal/util/validation"
)

// WatcherConfig holds raw data from the YAML config file.
type WatcherConfig struct {
	Name  string   `yaml:"name"`
	Nodes []string `yaml:"nodes"`
}

// Validate data from the YAML file.
func (s *WatcherConfig) Validate() {
	validationUtil.YAMLStruct(*s, nil)
}

// GetWatcherConfigs loads configs from the synchs directory.
func GetWatcherConfigs() []Config {
	fmt.Println("Watchers:")
	return ImportYAMLDir("./config/watch")
}
