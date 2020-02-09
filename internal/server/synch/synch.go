package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/lang"
	"github.com/christoph-karpowicz/unifier/internal/server/sim"
)

type synch struct {
	synch      *synchData
	dbs        map[string]*db.Database
	tables     map[string]*table
	nodes      map[string]*node
	mappings   []*mapping
	running    bool
	initial    bool
	simulation *sim.Simulation
}

func (s *synch) GetData() *synchData {
	return s.synch
}

func (s *synch) Init(DBMap map[string]*db.Database, simulation bool) {
	s.dbs = make(map[string]*db.Database)
	s.tables = make(map[string]*table)
	s.nodes = make(map[string]*node)

	if simulation {
		s.simulation = sim.CreateSimulation()
	}

	s.setDatabases(DBMap)
	s.setTables()
	s.setNodes()
	s.parseMappings()
	s.selectData()
	s.pairData()
}

// pairData pairs together records that are going to be synchronized.
func (s *synch) pairData() {
	for i := range s.mappings {
		var mpng *mapping = s.mappings[i]
		mpng.createPairs()
	}
}

func (s *synch) parseMappings() {
	for _, mapping := range s.synch.Mappings {
		rawMapping := lang.ParseMapping(mapping)
		for _, link := range rawMapping["links"].([]map[string]string) {
			mpng := createMapping(s, link, rawMapping["matchMethod"].(map[string]interface{}), rawMapping["do"].([]string))
			s.mappings = append(s.mappings, mpng)
		}
	}
}

// // Selects all records from all tables and filters them to get the relevant records.
func (s *synch) selectData() {
	for i := range s.mappings {
		var mpng *mapping = s.mappings[i]
		sourceRawActiveRecords := (*mpng.source.db).Select(mpng.source.tbl.name, mpng.sourceWhere)
		targetRawActiveRecords := (*mpng.target.db).Select(mpng.target.tbl.name, mpng.targetWhere)

		if !s.initial {
			mpng.sourceOldActiveRecords = mpng.sourceActiveRecords
			mpng.targetOldActiveRecords = mpng.targetActiveRecords
		}

		for _, sourceRecord := range sourceRawActiveRecords {
			sourceRecordPointer := mpng.source.tbl.records.FindRecordPointer(sourceRecord)
			mpng.sourceActiveRecords = append(mpng.sourceActiveRecords, sourceRecordPointer)
			sourceRecordPointer.ActiveIn = append(sourceRecordPointer.ActiveIn, mpng)
		}
		for _, targetRecord := range targetRawActiveRecords {
			targetRecordPointer := mpng.target.tbl.records.FindRecordPointer(targetRecord)
			mpng.targetActiveRecords = append(mpng.targetActiveRecords, targetRecordPointer)
			targetRecordPointer.ActiveIn = append(targetRecordPointer.ActiveIn, mpng)
		}
		// log.Println(mpng.sourceActiveRecords)
		// log.Println(mpng.targetActiveRecords)
	}
}

func (s *synch) setDatabase(DBMap map[string]*db.Database, dbName string) {
	_, dbExists := DBMap[dbName]
	if dbExists {
		s.dbs[dbName] = DBMap[dbName]
		(*s.dbs[dbName]).Init()
	} else {
		panic("[set database] ERROR: database " + dbName + " hasn't been configured.")
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
	for i := range s.synch.Nodes {
		var nodeData *nodeData = &s.synch.Nodes[i]

		var tableName string = nodeData.Database + "." + nodeData.Table
		_, tableFound := s.tables[tableName]
		if !tableFound {
			log.Fatalln("[create node] ERROR: table " + tableName + " not found.")
		}

		s.nodes[nodeData.Name] = createNode(nodeData, s.dbs[nodeData.Database], s.tables[tableName])
	}
}

// setTable creates an individual table struct and selects all records from it.
func (s *synch) setTable(tableName string, database *db.Database, key string) {
	var tblID string = (*database).GetData().GetName() + "." + tableName
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

func (s *synch) SynchPairs() {
	for j := range s.mappings {
		var mpng *mapping = s.mappings[j]

		for k := range mpng.pairs {
			var pair *pair = mpng.pairs[k]
			_, err := pair.synchronize(mpng.source.db, mpng.target.db)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
