package synch

import (
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

func (s *Synch) SelectData() {

	// Select all records from all tables.
	for _, table := range s.synch.Tables {
		rawDb1Records := (*s.database1).Select(table.Names.Table1, "")
		rawDb2Records := (*s.database2).Select(table.Names.Table2, "")
		table.Db1Records = TableRecords{records: MapToRecords(rawDb1Records)}
		table.Db2Records = TableRecords{records: MapToRecords(rawDb2Records)}

		for _, vector := range table.Vectors {
			rawDb1ActiveRecords := (*s.database1).Select(table.Names.Table1, vector.Conditions.Table1)
			rawDb2ActiveRecords := (*s.database2).Select(table.Names.Table2, vector.Conditions.Table2)
			for _, db1record := range rawDb1ActiveRecords {
				db1recordPointer := table.Db1Records.FindRecordPointer(db1record)
				vector.Db1ActiveRecords = append(vector.Db1ActiveRecords, db1recordPointer)
			}
			for _, db2record := range rawDb2ActiveRecords {
				db2recordPointer := table.Db2Records.FindRecordPointer(db2record)
				vector.Db2ActiveRecords = append(vector.Db2ActiveRecords, db2recordPointer)
				log.Println(*db2recordPointer)
			}
			log.Println(vector.Db1ActiveRecords)
			// 	vector.Db1ActiveRecords = MapToRecords(rawDb1ActiveRecords)
			// 	vector.Db2ActiveRecords = MapToRecords(rawDb2ActiveRecords)
		}
	}

	// log.Println(s.synch.Tables[0].Vectors[0])

}
