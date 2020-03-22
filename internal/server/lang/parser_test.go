package lang

import (
	"log"
	"testing"
)

var exampleInstruction string = `
MAP
	dvdrental_films.film_id TO msamp_films.ext_id,
	dvdrental_films.title TO msamp_films.Title,
	dvdrental_films.description TO msamp_films.Title,
	dvdrental_films.rental_duration TO msamp_films."Rental Duration",
	dvdrental_films.length TO msamp_films.Length,
	dvdrental_films.replacement_cost TO msamp_films."Replacement Cost",
	dvdrental_films.rating TO msamp_films.Rating,
	dvdrental_films.special_features TO msamp_films."Special Features"
SYNCH
	[dvdrental_films.title WHERE film_id <= 3] TO [msamp_films.Title],
	[dvdrental_films.title WHERE film_id > 30 AND film_id <= 50] TO [msamp_films.Title]
MATCH BY IDS(dvdrental_films.film_id, msamp_films.ext_id)
DO UPDATE, INSERT`

func TestParser(t *testing.T) {
	rawInstruction, err := ParseInstruction(exampleInstruction)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(rawInstruction)

	if rawInstruction["mapCmd"].(string) != "MAP" {
		log.Fatal("MAP command hasn't been read properly.")
	}
	mappingsLen := len(rawInstruction["mappings"].([]map[string]string))
	if mappingsLen != 8 {
		log.Fatalf("There should be 8 mappings, are %d.", mappingsLen)
	}
	if rawInstruction["synchCmd"].(string) != "SYNCH" {
		log.Fatal("SYNCH command hasn't been read properly.")
	}
	linksLen := len(rawInstruction["links"].([]map[string]string))
	if linksLen != 2 {
		log.Fatalf("There should be 2 links, are %d.", linksLen)
	}
	if rawInstruction["matchMethod"].(map[string]interface{})["matchCmd"] != "IDS" {
		log.Fatal("matchMethod hasn't been read properly.")
	}
	if rawInstruction["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[0] != "dvdrental_films.film_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if rawInstruction["matchMethod"].(map[string]interface{})["matchArgs"].([]string)[1] != "msamp_films.ext_id" {
		log.Fatal("matchArgs haven't been read properly.")
	}
	if len(rawInstruction["do"].([]string)) != 2 || (rawInstruction["do"].([]string)[0] != "UPDATE" || rawInstruction["do"].([]string)[1] != "INSERT") {
		log.Fatal("'do' action hasn't been read properly.")
	}

}
