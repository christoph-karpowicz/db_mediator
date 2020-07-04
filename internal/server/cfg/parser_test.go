package cfg

import (
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	os.Chdir("../../..")

	var synchCfgs []*SynchConfig = GetSynchConfigs()
	for _, config := range synchCfgs {
		if config.Name == "films" {
			config.Validate()

			for _, mapping := range config.Map {
				ParseMapping(mapping)
			}
			for _, link := range config.Link {
				ParseLink(link)
			}
		}
	}
}
