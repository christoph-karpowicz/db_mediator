package cfg

import (
	"log"
	"os"
	"testing"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
)

func TestParser(t *testing.T) {
	os.Chdir("../../..")

	var synchCfgs []Config = GetSynchConfigs()
	for _, config := range synchCfgs {
		cfg := config.(*cfg.SynchConfig)
		if cfg.Name == "films" {
			cfg.Validate()

			for _, mapping := range cfg.Map {
				_, err := ParseMapping(mapping)
				if err != nil {
					log.Fatalln(err)
				}
			}
			for _, link := range cfg.Link {
				_, err := ParseLink(link)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
