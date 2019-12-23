package synch

import (
	"fmt"
	"log"
)

type Vector struct {
	ColumnNames         TableSpecifics `json:"columnNames"`
	DataFlow            string         `json:"dataFlow"`
	Conditions          TableSpecifics `json:"conditions"`
	Db1OldActiveRecords []*Record
	Db2OldActiveRecords []*Record
	Db1ActiveRecords    []*Record
	Db2ActiveRecords    []*Record
	Pairs               []Pair
}

// For each active record in database1 find a corresponding acitve record in database2.
func (v *Vector) CreatePairs(settings Settings) {
	for i := range v.Db1ActiveRecords {
		DB1_record := v.Db1ActiveRecords[i]
		for j := range v.Db2ActiveRecords {
			DB2_record := v.Db2ActiveRecords[j]

			if settings.SynchType.MatchBy == "external_id_columns" {
				var DB1_externalIdColumnName string = settings.SynchType.ColumnNames.Table1
				var DB2_externalIdColumnName string = settings.SynchType.ColumnNames.Table2
				DB1_externalId, DB1_ok := DB1_record.Data[DB1_externalIdColumnName]
				DB2_externalId, DB2_ok := DB2_record.Data[DB2_externalIdColumnName]

				if !DB1_ok || !DB2_ok {
					continue
				}

				if areEqual, err := AreEqual(DB1_externalId, DB2_externalId); err != nil {
					log.Println(err)
				} else if areEqual {
					newPairs, err := CreatePair(DB1_record, DB2_record, v.DataFlow)
					if err != nil {
						log.Println(err)
					}

					for _, pair := range newPairs {
						v.Pairs = append(v.Pairs, pair)
					}
					DB1_record.PairedIn = append(DB1_record.PairedIn, v)
					DB2_record.PairedIn = append(DB2_record.PairedIn, v)
				}
			}

		}
	}
	for _, pair := range v.Pairs {
		fmt.Printf("rec1: %s\n", pair.GetSource().Data)
		fmt.Printf("rec2: %s\n", pair.GetTarget().Data)
		fmt.Println("======")
	}
}
