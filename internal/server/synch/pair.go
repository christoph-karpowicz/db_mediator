package synch

import (
	"fmt"
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/sim"
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
					src := &sim.RecordState{
						KeyName:      p.source.Key,
						KeyValue:     p.source.Data[p.source.Key],
						ColumnName:   p.mapping.sourceColumn,
						CurrentValue: p.source.Data[p.mapping.sourceColumn],
						NewValue:     nil,
					}
					trgt := &sim.RecordState{
						KeyName:      p.target.Key,
						KeyValue:     p.target.Data[p.target.Key],
						ColumnName:   p.mapping.targetColumn,
						CurrentValue: p.target.Data[p.mapping.targetColumn],
						NewValue:     sourceColumnValue,
					}

					p.mapping.synch.simulation.AddUpdate(src, trgt)
					fmt.Println(p.mapping.synch.simulation)
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
