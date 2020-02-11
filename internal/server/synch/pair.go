package synch

import (
	"fmt"
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type pair struct {
	mapping *mapping
	source  *record
	target  *record
}

func createPair(mpng *mapping, source *record, target *record) *pair {
	var newPair pair = pair{
		mapping: mpng,
		source:  source,
		target:  target,
	}

	return &newPair
}

func (p *pair) synchronize(db1 *db.Database, db2 *db.Database) (bool, error) {
	// Updates
	// If this pair is complete.
	if p.target != nil {
		// log.Println(p.source)
		// log.Println(p.target)

		if p.mapping.sourceColumn != "*" && p.mapping.targetColumn != "*" {
			sourceColumnValue := p.source.Data[p.mapping.sourceColumn]
			targetColumnValue := p.target.Data[p.mapping.targetColumn]

			if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
				log.Println(err)
			} else if !areEqual {
				if p.mapping.synch.simulation != nil {
					p.mapping.synch.simulation.AddUpdate(p)
					// fmt.Println(p.mapping.synch.simulation)
				}
				// (*db2).Update("", p.target.Key, p.mapping.targetColumn, sourceColumnValue)
				// log.Println(sourceColumnValue)
				// log.Println(targetColumnValue)
			}
		}
		// Inserts
		// If a target record has to be created.
	} else {

	}

	return false, nil
}

func (p *pair) CreateSimulationString() string {
	return fmt.Sprintf("|%6v: %6v, %6v: %16v| => |%6v: %6v, %6s: %16v -> %16v|\n",
		p.source.Key,
		p.source.Data[p.source.Key],
		p.mapping.sourceColumn,
		p.source.Data[p.mapping.sourceColumn],
		p.target.Key,
		p.target.Data[p.target.Key],
		p.mapping.targetColumn,
		p.target.Data[p.mapping.targetColumn],
		p.source.Data[p.mapping.sourceColumn],
	)
}
