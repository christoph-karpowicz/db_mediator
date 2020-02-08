package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type pair struct {
	mapping    *mapping
	source     *record
	target     *record
	IsComplete bool // does the source record have a corresponding target record
}

func (p *pair) synchronize(db1 *db.Database, db2 *db.Database) (bool, error) {
	// Updates
	if p.IsComplete {
		log.Println(p.source)
		log.Println(p.target)

		if p.mapping.sourceColumn != "*" && p.mapping.targetColumn != "*" {
			sourceColumnValue := p.source.Data[p.mapping.sourceColumn]
			targetColumnValue := p.target.Data[p.mapping.targetColumn]

			if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
				log.Println(err)
			} else if !areEqual {
				// (*db2).Update("", p.target.Key, p.mapping.targetColumn, sourceColumnValue)
				// log.Println(sourceColumnValue)
				// log.Println(targetColumnValue)
			}
		}
		// Inserts
	} else {

	}

	// if secondaryFlow != nil {

	// }
	return false, nil
}

func createPair(mpng *mapping, source *record, target *record) pair {
	var newPair pair = pair{
		mapping: mpng,
		source:  source,
		target:  target,
	}

	if target != nil {
		newPair.IsComplete = true
	} else {
		newPair.IsComplete = false
	}

	return newPair
}
