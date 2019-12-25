package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type synch struct {
	synch     *synchData
	database1 *db.Database
	database2 *db.Database
	running   bool
	initial   bool
}

func (s *synch) GetData() *synchData {
	return s.synch
}

func (s *synch) Init(DBMap map[string]*db.Database) {
	s.setDatabases(DBMap)
	s.selectData()
	s.pairData()
	s.setParentPointers()
}

// pairData pairs together records that are going to be synchronized.
func (s *synch) pairData() {
	for i := range s.synch.Tables {
		var table *table = &s.synch.Tables[i]

		for j := range table.Vectors {
			var vector *vector = &table.Vectors[j]
			vector.createPairs()
		}
	}
}

// Selects all records from all tables and filters them to get the relevant records.
func (s *synch) selectData() {
	for i := range s.synch.Tables {
		var table *table = &s.synch.Tables[i]
		DB1_rawRecords := (*s.database1).Select(table.Names.Table1, "-")
		DB2_rawRecords := (*s.database2).Select(table.Names.Table2, "-")

		if !s.initial {
			table.Db1OldRecords = table.Db1Records
			table.Db2OldRecords = table.Db2Records
		}

		table.Db1Records = &tableRecords{records: mapToRecords(DB1_rawRecords, table.Keys.Table1)}
		table.Db2Records = &tableRecords{records: mapToRecords(DB2_rawRecords, table.Keys.Table2)}

		for j := range table.Vectors {
			var vector *vector = &table.Vectors[j]
			DB1_rawActiveRecords := (*s.database1).Select(table.Names.Table1, vector.Conditions.Table1)
			DB2_rawActiveRecords := (*s.database2).Select(table.Names.Table2, vector.Conditions.Table2)

			if !s.initial {
				vector.Db1OldActiveRecords = vector.Db1ActiveRecords
				vector.Db2OldActiveRecords = vector.Db2ActiveRecords
			}

			for _, DB1_record := range DB1_rawActiveRecords {
				DB1_recordPointer := table.Db1Records.FindRecordPointer(DB1_record)
				vector.Db1ActiveRecords = append(vector.Db1ActiveRecords, DB1_recordPointer)
				DB1_recordPointer.ActiveIn = append(DB1_recordPointer.ActiveIn, vector)
			}
			for _, DB2_record := range DB2_rawActiveRecords {
				DB2_recordPointer := table.Db2Records.FindRecordPointer(DB2_record)
				vector.Db2ActiveRecords = append(vector.Db2ActiveRecords, DB2_recordPointer)
				DB2_recordPointer.ActiveIn = append(DB2_recordPointer.ActiveIn, vector)
			}
			log.Println(vector.Db1ActiveRecords)
			log.Println(vector.Db2ActiveRecords)
		}
	}
}

// Open chosen database connections.
func (s *synch) setDatabases(DBMap map[string]*db.Database) {
	if DBMap[s.synch.Databases.Db1.Name] == nil || DBMap[s.synch.Databases.Db2.Name] == nil {
		panic(s.synch.Name + " database config is invalid.")
	}
	s.database1 = DBMap[s.synch.Databases.Db1.Name]
	(*s.database1).Init()
	s.database2 = DBMap[s.synch.Databases.Db2.Name]
	(*s.database2).Init()
}

func (s *synch) setParentPointers() {
	for i := range s.synch.Tables {
		var table *table = &s.synch.Tables[i]

		for j := range table.Vectors {
			var vector *vector = &table.Vectors[j]
			vector.table = table

			for k := range vector.pairs {
				var pair *pair = &table.Vectors[j].pairs[k]
				pair.vector = vector
			}
		}
	}
}

func (s *synch) SynchPairs() {
	for i := range s.synch.Tables {
		var table *table = &s.synch.Tables[i]

		for j := range table.Vectors {
			var vector *vector = &table.Vectors[j]

			for k := range vector.pairs {
				var pair *pair = &vector.pairs[k]
				_, err := pair.synchronize(s.database1, s.database2)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
