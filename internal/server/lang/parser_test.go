package lang

import (
	"fmt"
	"log"
	"testing"
)

var exampleMapping string = `
MAP
	[dvdrental_films.title WHERE title ILIKE 'A%'] TO [msamp_films.Title], 
	[dvdrental_films.actor] TO [msamp_films.Actor WHERE {actor: "Daniel Day Lewis"}],
	[dvdrental_films.year] TO [msamp_films.Year]
MATCH BY IDS(dvdrental_films.film_id, msamp_films.ext_id)
DO UPDATE`

func TestParser(t *testing.T) {
	rawMapping := ParseMapping(exampleMapping)
	fmt.Println(rawMapping)

	if rawMapping["command"].(string) != "MAP" {
		log.Fatal("Command hasn't been read properly.")
	}
	if len(rawMapping["links"].([]map[string]string)) == 0 {
		log.Fatal("There should be 3 links.")
	}
	if rawMapping["matchMethod"].(map[string]interface{})["matchCmd"] != "IDS" {
		log.Fatal("matchMethod hasn't been read properly.")
	}
	if rawMapping["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[0] != "dvdrental_films.film_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if rawMapping["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[1] != "msamp_films.ext_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if len(rawMapping["do"].([]string)) == 1 && rawMapping["do"].([]string)[0] != "UPDATE" {
		log.Fatal("'do' action hasn't been read properly.")
	}

}