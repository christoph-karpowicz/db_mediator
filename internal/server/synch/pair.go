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

func (p pair) getSourceNodeKey() string {
	return p.mapping.source.data.Key
}

func (p pair) getTargetNodeKey() string {
	return p.mapping.target.data.Key
}

func (p pair) synchronize(db1 *db.Database, db2 *db.Database) (bool, error) {
	if p.target != nil {
		// Updates
		// If this pair is complete.
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
			} else {
				if p.mapping.synch.simulation != nil {
					p.mapping.synch.simulation.AddIdle(p)
				}
			}
		}
	} else {
		// Inserts
		// If a target record has to be created.
		if p.mapping.synch.simulation != nil {
			p.mapping.synch.simulation.AddInsert(p)
			// fmt.Println(p.mapping.synch.simulation)
		}
	}

	return false, nil
}

func (p pair) SimIdleString() string {
	return fmt.Sprintf("|%6v: %6v, %6v: %25v|  ==  |%6v: %6v, %6s: %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.mapping.sourceColumn,
		p.source.Data[p.mapping.sourceColumn],
		p.getTargetNodeKey(),
		p.target.Data[p.getTargetNodeKey()],
		p.mapping.targetColumn,
		p.target.Data[p.mapping.targetColumn],
	)
}

func (p pair) SimInsertString() string {
	return fmt.Sprintf("|%6v: %6v, %6v: %25v|  =>  |%6v: %6v, %6s: %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.mapping.sourceColumn,
		p.source.Data[p.mapping.sourceColumn],
		p.getTargetNodeKey(),
		"-",
		p.mapping.targetColumn,
		p.source.Data[p.mapping.sourceColumn],
	)
}

func (p pair) SimUpdateString() string {
	return fmt.Sprintf("|%6v: %6v, %6v: %25v|  =^  |%6v: %6v, %6s: %16v -> %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.mapping.sourceColumn,
		p.source.Data[p.mapping.sourceColumn],
		p.getTargetNodeKey(),
		p.target.Data[p.getTargetNodeKey()],
		p.mapping.targetColumn,
		p.target.Data[p.mapping.targetColumn],
		p.source.Data[p.mapping.sourceColumn],
	)
}
