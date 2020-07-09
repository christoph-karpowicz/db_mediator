package cfg

import (
	"log"
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
				_, err := ParseMapping(mapping)
				if err != nil {
					log.Fatalln(err)
				}
			}
			for _, link := range config.Link {
				_, err := ParseLink(link)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
