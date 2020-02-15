package synch

import (
	"fmt"
	"log"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

type Pair struct {
	mapping *mapping
	source  *record
	target  *record
}

func createPair(mpng *mapping, source *record, target *record) *Pair {
	var newPair Pair = Pair{
		mapping: mpng,
		source:  source,
		target:  target,
	}

	return &newPair
}

func (p Pair) getSourceNodeKey() string {
	return p.mapping.source.data.Key
}

func (p Pair) getTargetNodeKey() string {
	return p.mapping.target.data.Key
}

func (p Pair) Synchronize() (bool, error) {
	// db1 := p.mapping.source.db
	// db2 := p.mapping.target.db

	if p.target != nil && arrUtil.Contains(p.mapping.do, "UPDATE") {
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
				if p.mapping.synch.Simulation != nil {
					p.mapping.synch.Simulation.AddUpdate(p)
					// fmt.Println(p.mapping.synch.Simulation)
				}
				// (*db2).Update("", p.target.Key, p.mapping.targetColumn, sourceColumnValue)
				// log.Println(sourceColumnValue)
				// log.Println(targetColumnValue)
			} else {
				if p.mapping.synch.Simulation != nil {
					p.mapping.synch.Simulation.AddIdle(p)
				}
			}
		}
	} else if p.target == nil && arrUtil.Contains(p.mapping.do, "INSERT") {
		// Inserts
		// If a target record has to be created.
		if p.mapping.synch.Simulation != nil {
			p.mapping.synch.Simulation.AddInsert(p)
			// fmt.Println(p.mapping.synch.Simulation)
		}
	}

	return false, nil
}

func (p Pair) SimIdleString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  ==  |%6v: %6v, %6s: %25v|\n",
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

func (p Pair) SimInsertString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =>  |%6v: %6v, %6s: %25v|\n",
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

func (p Pair) SimUpdateString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =^  |%6v: %6v, %6s: %25v -> %25v|\n",
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
