package synch

import (
	"fmt"
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type Synch struct {
	synch     *SynchData
	database1 *db.Database
	database2 *db.Database
}

func (s *Synch) GetData() *SynchData {
	return s.synch
}

func (s *Synch) SetDatabases(DBMap map[string]*db.Database) {
	if DBMap[s.synch.Databases.Db1.Name] == nil || DBMap[s.synch.Databases.Db2.Name] == nil {
		panic(s.synch.Name + " database config is invalid.")
	}
	s.database1 = DBMap[s.synch.Databases.Db1.Name]
	(*s.database1).Init()
	s.database2 = DBMap[s.synch.Databases.Db2.Name]
	(*s.database2).Init()
}

func (s *Synch) PairData() {
	for _, table := range s.synch.Tables {
		for _, vector := range table.Vectors {
			// For each active record in database1 find a corresponding acitve record in database2.
			fmt.Println(table.Names.Table1)
			fmt.Println(table.Names.Table2)
			fmt.Println(table.Db1Records)
			for _, DB1_record := range vector.Db1ActiveRecords {
				for _, DB2_record := range vector.Db2ActiveRecords {

					if table.Settings.SynchType.MatchBy == "external_id_columns" {
						var DB1_externalIdColumnName string = table.Settings.SynchType.ColumnNames.Table1
						var DB2_externalIdColumnName string = table.Settings.SynchType.ColumnNames.Table2
						if DB1_record.Data[DB1_externalIdColumnName] == DB2_record.Data[DB2_externalIdColumnName] {
							var newPair Pair = Pair{record1: DB1_record, record2: DB2_record}
							vector.Pairs = append(vector.Pairs, newPair)
						}
					}

				}
			}
			for _, pair := range vector.Pairs {
				fmt.Printf("rec1: %s, rec2: %s\n", pair.record1.Data, pair.record2.Data)
				fmt.Println("======")
			}
		}
	}
}

func (s *Synch) SelectData() {

	// Select all records from all tables.
	for _, table := range s.synch.Tables {
		DB1_rawRecords := (*s.database1).Select(table.Names.Table1, "")
		DB2_rawRecords := (*s.database2).Select(table.Names.Table2, "")
		table.Db1Records = &TableRecords{records: MapToRecords(DB1_rawRecords)}
		table.Db2Records = &TableRecords{records: MapToRecords(DB2_rawRecords)}

		for _, vector := range table.Vectors {
			DB1_rawActiveRecords := (*s.database1).Select(table.Names.Table1, vector.Conditions.Table1)
			DB2_rawActiveRecords := (*s.database2).Select(table.Names.Table2, vector.Conditions.Table2)
			for _, DB1_record := range DB1_rawActiveRecords {
				DB1_recordPointer := table.Db1Records.FindRecordPointer(DB1_record)
				vector.Db1ActiveRecords = append(vector.Db1ActiveRecords, DB1_recordPointer)
			}
			for _, DB2_record := range DB2_rawActiveRecords {
				DB2_recordPointer := table.Db2Records.FindRecordPointer(DB2_record)
				vector.Db2ActiveRecords = append(vector.Db2ActiveRecords, DB2_recordPointer)
				// log.Println(*db2recordPointer)
			}
			log.Println(vector.Db1ActiveRecords)
			log.Println(vector.Db2ActiveRecords)
		}
	}

	// log.Println(s.synch.Tables[0].Vectors[0])

}
