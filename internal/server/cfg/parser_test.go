package cfg

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	rawLink, err := ParseLink(exampleLink)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(rawLink)

	if rawLink["mapCmd"].(string) != "MAP" {
		log.Fatal("MAP command hasn't been read properly.")
	}
	mappingsLen := len(rawLink["mappings"].([]map[string]string))
	if mappingsLen != 8 {
		log.Fatalf("There should be 8 mappings, are %d.", mappingsLen)
	}
	if rawLink["synchCmd"].(string) != "SYNCH" {
		log.Fatal("SYNCH command hasn't been read properly.")
	}
	linksLen := len(rawLink["links"].([]map[string]string))
	if linksLen != 2 {
		log.Fatalf("There should be 2 links, are %d.", linksLen)
	}
	if rawLink["matchMethod"].(map[string]interface{})["matchCmd"] != "IDS" {
		log.Fatal("matchMethod hasn't been read properly.")
	}
	if rawLink["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[0] != "dvdrental_films.film_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if rawLink["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[1] != "msamp_films.ext_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if len(rawLink["do"].([]string)) != 2 || (rawLink["do"].([]string)[0] != "UPDATE" || rawLink["do"].([]string)[1] != "INSERT") {
		log.Fatal("'do' action hasn't been read properly.")
	}

}