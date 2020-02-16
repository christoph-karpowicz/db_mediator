package synch

import (
	"fmt"
	"log"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

type Pair struct {
	Mapping *Mapping
	source  *record
	target  *record
}

func createPair(mpng *Mapping, source *record, target *record) *Pair {
	var newPair Pair = Pair{
		Mapping: mpng,
		source:  source,
		target:  target,
	}

	return &newPair
}

func (p Pair) getSourceNodeKey() string {
	return p.Mapping.source.data.Key
}

func (p Pair) getTargetNodeKey() string {
	return p.Mapping.target.data.Key
}

func (p Pair) Synchronize() (bool, error) {
	// db1 := p.Mapping.source.db
	// db2 := p.Mapping.target.db

	if p.target != nil && arrUtil.Contains(p.Mapping.do, "UPDATE") {
		// Updates
		// If this pair is complete.
		// log.Println(p.source)
		// log.Println(p.target)

		if p.Mapping.sourceColumn != "*" && p.Mapping.targetColumn != "*" {
			sourceColumnValue := p.source.Data[p.Mapping.sourceColumn]
			targetColumnValue := p.target.Data[p.Mapping.targetColumn]

			if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
				log.Println(err)
			} else if !areEqual {
				if p.Mapping.synch.Simulation != nil {
					p.Mapping.synch.Simulation.AddUpdate(p)
					// fmt.Println(p.Mapping.synch.Simulation)
				}
				// (*db2).Update("", p.target.Key, p.Mapping.targetColumn, sourceColumnValue)
				// log.Println(sourceColumnValue)
				// log.Println(targetColumnValue)
			} else {
				if p.Mapping.synch.Simulation != nil {
					p.Mapping.synch.Simulation.AddIdle(p)
				}
			}
		}
	} else if p.target == nil && arrUtil.Contains(p.Mapping.do, "INSERT") {
		// Inserts
		// If a target record has to be created.
		if p.Mapping.synch.Simulation != nil {
			p.Mapping.synch.Simulation.AddInsert(p)
			// fmt.Println(p.Mapping.synch.Simulation)
		}
	}

	return false, nil
}

func (p Pair) SimIdleString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  ==  |%6v: %6v, %6s: %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.Mapping.sourceColumn,
		p.source.Data[p.Mapping.sourceColumn],
		p.getTargetNodeKey(),
		p.target.Data[p.getTargetNodeKey()],
		p.Mapping.targetColumn,
		p.target.Data[p.Mapping.targetColumn],
	)
}

func (p Pair) SimInsertString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =>  |%6v: %6v, %6s: %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.Mapping.sourceColumn,
		p.source.Data[p.Mapping.sourceColumn],
		p.getTargetNodeKey(),
		"-",
		p.Mapping.targetColumn,
		p.source.Data[p.Mapping.sourceColumn],
	)
}

func (p Pair) SimUpdateString() string {
	return fmt.Sprintf("|%6v: %3v, %6v: %25v|  =^  |%6v: %6v, %6s: %25v -> %25v|\n",
		p.getSourceNodeKey(),
		p.source.Data[p.getSourceNodeKey()],
		p.Mapping.sourceColumn,
		p.source.Data[p.Mapping.sourceColumn],
		p.getTargetNodeKey(),
		p.target.Data[p.getTargetNodeKey()],
		p.Mapping.targetColumn,
		p.target.Data[p.Mapping.targetColumn],
		p.source.Data[p.Mapping.sourceColumn],
	)
}
