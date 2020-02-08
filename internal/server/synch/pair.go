package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type pair struct {
	mapping    *mapping
	source     *record
	target     *record
	IsComplete bool
}

func (p *pair) synchronize(db1 *db.Database, db2 *db.Database) (bool, error) {
	log.Println(p.source)
	log.Println(p.target)
	// // Updates
	// if p.IsComplete {
	// 	if p.primaryFlow.sourceColumnName != "*" && p.primaryFlow.targetColumnName != "*" {
	// 		source := p.primaryFlow.source
	// 		target := p.primaryFlow.target
	// 		sourceColumnValue := source.Data[p.primaryFlow.sourceColumnName]
	// 		targetColumnValue := target.Data[p.primaryFlow.targetColumnName]

	// 		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
	// 			log.Println(err)
	// 		} else if !areEqual {
	// 			(*db2).Update("", target.Key, p.primaryFlow.targetColumnName, sourceColumnValue)
	// 			// log.Println(sourceColumnValue)
	// 			// log.Println(targetColumnValue)
	// 		}
	// 	}
	// 	// Inserts
	// } else {

	// }

	// // if secondaryFlow != nil {

	// // }
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
