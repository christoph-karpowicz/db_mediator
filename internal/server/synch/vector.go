package synch

import (
	"fmt"
	"log"
)

type vector struct {
	ColumnNames         tableSpecifics `json:"columnNames"`
	DataFlow            string         `json:"dataFlow"`
	Conditions          tableSpecifics `json:"conditions"`
	Db1OldActiveRecords []*record
	Db2OldActiveRecords []*record
	Db1ActiveRecords    []*record
	Db2ActiveRecords    []*record
	Pairs               []pair
}

// For each active record in database1 find a corresponding acitve record in database2.
func (v *vector) createPairs(settings settings) {
	var sourceRecords []*record
	var targetRecords []*record
	var isBidirectional bool = false

	if v.DataFlow == "=>" || v.DataFlow == "<=>*" {
		sourceRecords = v.Db1ActiveRecords
		targetRecords = v.Db2ActiveRecords
	} else {
		sourceRecords = v.Db2ActiveRecords
		targetRecords = v.Db1ActiveRecords
	}

	if v.DataFlow == "*<=>" || v.DataFlow == "<=>*" {
		isBidirectional = true
	}

	for i := range sourceRecords {
		source := sourceRecords[i]
		var pairFound bool = false
		for j := range targetRecords {
			target := targetRecords[j]

			if settings.SynchType.MatchBy == "external_id_columns" {
				var sourceExternalIdColumnName string = settings.SynchType.ColumnNames.Table1
				var targetExternalIdColumnName string = settings.SynchType.ColumnNames.Table2
				sourceExternalId, sourceOk := source.Data[sourceExternalIdColumnName]
				targetExternalId, targetOk := target.Data[targetExternalIdColumnName]

				if !sourceOk || !targetOk {
					continue
				}

				if areEqual, err := areEqual(sourceExternalId, targetExternalId); err != nil {
					log.Println(err)
				} else if areEqual {
					newPair := createPair(source, target, v.DataFlow, v.ColumnNames)
					v.Pairs = append(v.Pairs, newPair)
					pairFound = true
					source.PairedIn = append(source.PairedIn, v)
					target.PairedIn = append(target.PairedIn, v)
				}
			}
		}
		if !pairFound && isBidirectional {
			newPair := createPair(source, nil, v.DataFlow, v.ColumnNames)
			v.Pairs = append(v.Pairs, newPair)
		}
	}
	for _, pair := range v.Pairs {
		fmt.Printf("rec1: %s\n", pair.primaryFlow.source.Data)
		if pair.IsComplete {
			fmt.Printf("rec2: %s\n", pair.primaryFlow.target.Data)
		}
		fmt.Println("======")
	}
}
