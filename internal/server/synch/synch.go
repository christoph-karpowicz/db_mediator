package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/lang"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// Synch represents an individual synchronzation configration.
// It holds all config data from an .yaml file, raw and parsed.
type Synch struct {
	Cfg        *Config
	dbs        map[string]*db.Database
	tables     map[string]*table
	nodes      map[string]*node
	Mappings   []*Mapping
	running    bool
	initial    bool
	Simulation unifier.Simulator
}

// GetConfig returns raw synch data.
func (s *Synch) GetConfig() *Config {
	return s.Cfg
}

// Init prepares the synchronization by fetching all necessary data
// and parsing it.
func (s *Synch) Init(DBMap map[string]*db.Database, simulation bool) {
	s.dbs = make(map[string]*db.Database)
	s.tables = make(map[string]*table)
	s.nodes = make(map[string]*node)

	s.setDatabases(DBMap)
	s.setTables()
	s.setNodes()
	s.parseMappings()
	s.selectData()
	s.pairData()
}

// pairData pairs together records that are going to be synchronized.
func (s *Synch) pairData() {
	for i := range s.Mappings {
		var mpng *Mapping = s.Mappings[i]
		mpng.createPairs()
	}
}

func (s *Synch) parseMappings() {
	for i, mapping := range s.Cfg.Mappings {
		rawMapping := lang.ParseMapping(mapping)
		for j, link := range rawMapping["links"].([]map[string]string) {
			mpng := createMapping(s, link, rawMapping["matchMethod"].(map[string]interface{}), rawMapping["do"].([]string), i, j)
			s.Mappings = append(s.Mappings, mpng)
		}
	}
}

// selectData selects all records from all tables and filters them to get the relevant records.
func (s *Synch) selectData() {
	for i := range s.Mappings {
		var mpng *Mapping = s.Mappings[i]
		sourceRawActiveRecords := (*mpng.source.db).Select(mpng.source.tbl.name, mpng.sourceWhere)
		targetRawActiveRecords := (*mpng.target.db).Select(mpng.target.tbl.name, mpng.targetWhere)

		// for _, v := range sourceRawActiveRecords {
		// 	fmt.Println(v["film_id"])
		// }

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

func (s *Synch) setDatabase(DBMap map[string]*db.Database, dbName string) {
	_, dbExists := DBMap[dbName]
	if dbExists {
		s.dbs[dbName] = DBMap[dbName]
		(*s.dbs[dbName]).Init()
	} else {
		panic("[set database] ERROR: database " + dbName + " hasn't been configured.")
	}
}

// setDatabases opens the chosen database connections.
func (s *Synch) setDatabases(DBMap map[string]*db.Database) {
	for j := range s.Cfg.Nodes {
		var nodeConfig *nodeConfig = &s.Cfg.Nodes[j]
		s.setDatabase(DBMap, nodeConfig.Database)
	}
}

// setNodes creates node structs and adds them to the relevant synch struct field.
func (s *Synch) setNodes() {
	for i := range s.Cfg.Nodes {
		var nodeConfig *nodeConfig = &s.Cfg.Nodes[i]

		var tableName string = nodeConfig.Database + "." + nodeConfig.Table
		_, tableFound := s.tables[tableName]
		if !tableFound {
			log.Fatalln("[create node] ERROR: table " + tableName + " not found.")
		}

		s.nodes[nodeConfig.Name] = createNode(nodeConfig, s.dbs[nodeConfig.Database], s.tables[tableName])
	}
}

// setTable creates an individual table struct and selects all records from it.
func (s *Synch) setTable(tableName string, database *db.Database) {
	var tblID string = (*database).GetConfig().GetName() + "." + tableName
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

		tbl.records = &tableRecords{records: mapToRecords(rawRecords)}
		s.tables[tbl.id] = tbl
	}
}

// setTables creates table structs based on node yaml data.
func (s *Synch) setTables() {
	for j := range s.Cfg.Nodes {
		var nodeConfig *nodeConfig = &s.Cfg.Nodes[j]
		s.setTable(nodeConfig.Table, s.dbs[nodeConfig.Database])
	}
}

// Synchronize loops over all pairs in all mappings and invokes their Synchronize function.
func (s *Synch) Synchronize() ([]byte, error) {
	for j := range s.Mappings {
		var mpng *Mapping = s.Mappings[j]

		for k := range mpng.pairs {
			var pair *Pair = mpng.pairs[k]
			_, err := pair.Synchronize()
			if err != nil {
				log.Println(err)
			}
		}
	}

	if s.Simulation != nil {
		return s.Simulation.ToJSON()
	}

	return []byte{}, nil
}
