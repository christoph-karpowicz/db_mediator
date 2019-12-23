package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type Pair struct {
	primaryFlow   Flow
	secondaryFlow Flow
	IsComplete    bool
}

func (p *Pair) Synchronize(db1 *db.Database, db2 *db.Database) (bool, error) {
	// Updates
	if p.IsComplete {
		if p.primaryFlow.sourceColumnName != "*" && p.primaryFlow.targetColumnName != "*" {
			source := p.primaryFlow.source
			target := p.primaryFlow.target
			sourceColumnValue := source.Data[p.primaryFlow.sourceColumnName]
			targetColumnValue := target.Data[p.primaryFlow.targetColumnName]

			if areEqual, err := AreEqual(sourceColumnValue, targetColumnValue); err != nil {
				log.Println(err)
			} else if !areEqual {
				(*db2).Update(target.Key, sourceColumnValue)
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

func CreatePair(source *Record, target *Record, flowSymbol string, columnNames TableSpecifics) Pair {
	var newPair Pair

	newPair.primaryFlow = Flow{
		source:           source,
		target:           target,
		sourceColumnName: columnNames.Table1,
		targetColumnName: columnNames.Table2,
	}

	if flowSymbol == "*<=>" || flowSymbol == "<=>*" {
		newPair.secondaryFlow = Flow{
			source:           source,
			target:           target,
			sourceColumnName: columnNames.Table1,
			targetColumnName: columnNames.Table2,
		}
	}

	if target != nil {
		newPair.IsComplete = true
	} else {
		newPair.IsComplete = false
	}

	return newPair
}
