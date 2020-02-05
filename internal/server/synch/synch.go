package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/lang"
)

type synch struct {
	synch    *synchData
	dbs      map[string]*db.Database
	tables   map[string]*table
	nodes    map[string]*node
	mappings []*mapping
	running  bool
	initial  bool
}

func (s *synch) GetData() *synchData {
	return s.synch
}

func (s *synch) Init(DBMap map[string]*db.Database) {
	s.dbs = make(map[string]*db.Database)
	s.tables = make(map[string]*table)
	s.setDatabases(DBMap)
	s.setTables()
	s.setNodes()
	// s.copyTables()
	// s.assignTablePointers()
	// s.selectData()
	// s.pairData()
	// s.setParentPointers()
}

// func (s *synch) assignTablePointers() {
// 	for i := range s.synch.Vectors {
// 		var vctr *vector = &s.synch.Vectors[i]

// 		vctr.sourceTable = s.tables[vctr.Source.Database+"."+vctr.Source.Table]
// 		vctr.targetTable = s.tables[vctr.Target.Database+"."+vctr.Target.Table]
// 	}
// }

// // pairData pairs together records that are going to be synchronized.
// func (s *synch) pairData() {
// 	for i := range s.synch.Vectors {
// 		var vector *vector = &s.synch.Vectors[i]
// 		vector.createPairs()
// 	}
// }

func (s *synch) parseMappings() {
	for _, mapping := range s.synch.Mappings {
		rawMapping := lang.ParseMapping(mapping)
		for _, link := range rawMapping["links"].([]map[string]string) {
			parsedMapping := createMapping(s.nodes, link, rawMapping["matchMethod"].(map[string]interface{}), rawMapping["do"].([]string))
			s.mappings = append(s.mappings, parsedMapping)
		}
	}
}

// // Selects all records from all tables and filters them to get the relevant records.
// func (s *synch) selectData() {
// 	for i := range s.synch.Vectors {
// 		var vctr *vector = &s.synch.Vectors[i]
// 		sourceRawActiveRecords := (*s.dbs[vctr.Source.Database]).Select(vctr.Source.Table, vctr.Source.Condition)
// 		targetRawActiveRecords := (*s.dbs[vctr.Target.Database]).Select(vctr.Target.Table, vctr.Target.Condition)

// 		if !s.initial {
// 			vctr.sourceOldActiveRecords = vctr.sourceActiveRecords
// 			vctr.targetOldActiveRecords = vctr.targetActiveRecords
// 		}

// 		for _, sourceRecord := range sourceRawActiveRecords {
// 			sourceRecordPointer := vctr.sourceTable.records.FindRecordPointer(sourceRecord)
// 			vctr.sourceActiveRecords = append(vctr.sourceActiveRecords, sourceRecordPointer)
// 			sourceRecordPointer.ActiveIn = append(sourceRecordPointer.ActiveIn, vctr)
// 		}
// 		for _, targetRecord := range targetRawActiveRecords {
// 			targetRecordPointer := vctr.targetTable.records.FindRecordPointer(targetRecord)
// 			vctr.targetActiveRecords = append(vctr.targetActiveRecords, targetRecordPointer)
// 			targetRecordPointer.ActiveIn = append(targetRecordPointer.ActiveIn, vctr)
// 		}
// 		// log.Println(vctr.sourceActiveRecords)
// 		// log.Println(vctr.targetActiveRecords)
// 	}
// }

func (s *synch) setDatabase(DBMap map[string]*db.Database, dbName string) {
	_, dbExists := DBMap[dbName]
	if dbExists {
		s.dbs[dbName] = DBMap[dbName]
		(*s.dbs[dbName]).Init()
	} else {
		panic("Database " + dbName + " hasn't been configured.")
	}
}

// Open chosen database connections.
func (s *synch) setDatabases(DBMap map[string]*db.Database) {
	for j := range s.synch.Nodes {
		var nodeData *nodeData = &s.synch.Nodes[j]
		s.setDatabase(DBMap, nodeData.Database)
	}
}

// setNodes creates node structs and adds them to the relevant synch struct field.
func (s *synch) setNodes() {
	for j := range s.synch.Nodes {
		var nodeData *nodeData = &s.synch.Nodes[j]
		s.nodes[nodeData.Name] = createNode(nodeData, s.dbs[nodeData.Database], s.tables[nodeData.Table])
	}
}

// func (s *synch) setParentPointers() {
// 	for j := range s.synch.Vectors {
// 		var vector *vector = &s.synch.Vectors[j]
// 		// vector.table = table

// 		for k := range vector.pairs {
// 			var pair *pair = &s.synch.Vectors[j].pairs[k]
// 			pair.vector = vector
// 		}
// 	}
// }

// setTable creates an individual table struct and selects all records from it.
func (s *synch) setTable(tableName string, database *db.Database, key string) {
	var tblID string = (*database).GetData().Name + "." + tableName
	_, tableCopied := s.tables[tblID]

	if !tableCopied {
		tbl := &table{
			id:   tblID,
			db:   database,
			name: tableName,
		}
		rawRecords := (*tbl.db).Select(tbl.name, "")

		if !s.initial {
			tbl.oldRecords = tbl.records
		}

		tbl.records = &tableRecords{records: mapToRecords(rawRecords, key)}
		s.tables[tbl.id] = tbl
	}
}

// setTables creates table structs based on node yaml data.
func (s *synch) setTables() {
	for j := range s.synch.Nodes {
		var nodeData *nodeData = &s.synch.Nodes[j]
		s.setTable(nodeData.Table, s.dbs[nodeData.Database], nodeData.Key)
	}
}

// func (s *synch) SynchPairs() {
// 	for j := range s.synch.Vectors {
// 		var vctr *vector = &s.synch.Vectors[j]

// 		for k := range vctr.pairs {
// 			var pair *pair = &vctr.pairs[k]
// 			_, err := pair.synchronize(s.dbs[vctr.Source.Database], s.dbs[vctr.Target.Database])
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
// 	}
// }
