/*
Package synch handles all data sychronization.
*/
package synch

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	"github.com/christoph-karpowicz/unifier/internal/server/lang"
	"github.com/christoph-karpowicz/unifier/internal/server/unifier"
)

// Synch represents an individual synchronzation configration.
// It holds all configuration from an .yaml file, raw and parsed.
type Synch struct {
	Cfg        *Config
	dbs        map[string]*db.Database
	tables     map[string]*table
	nodes      map[string]*node
	mappings   []*Mapping
	Links      []*Link
	running    bool
	initial    bool
	Simulation bool
	Rep        unifier.Reporter
}

// GetConfig returns the synch config struct.
func (s *Synch) GetConfig() *Config {
	return s.Cfg
}

// Init prepares the synchronization by fetching all necessary data
// and parsing it.
func (s *Synch) Init(DBMap map[string]*db.Database) {
	tStart := time.Now()
	s.dbs = make(map[string]*db.Database)
	s.tables = make(map[string]*table)
	s.nodes = make(map[string]*node)

	s.setDatabases(DBMap)
	s.setTables()
	s.setNodes()
	s.parseLinks()
	s.parseMappings()
	s.selectData()
	s.pairData()

	// fmt.Println(runtime.NumCPU())
	fmt.Println("Synch init finished in: ", time.Since(tStart).String())
}

// pairData pairs together records that are going to be synchronized.
func (s *Synch) pairData() {
	var wg sync.WaitGroup

	for i := range s.Links {
		var lnk *Link = s.Links[i]

		wg.Add(1)
		go lnk.createPairs(&wg)
		wg.Wait()
	}

}

func (s *Synch) parseLink(mpngStr string, i int, c chan bool) {
	rawLink, err := lang.ParseLink(mpngStr)
	if err != nil {
		panic(err)
	}

	fmt.Println(rawLink)
	in := createLink(s, rawLink)
	s.Links = append(s.Links, in)

	c <- true
}

func (s *Synch) parseLinks() {
	var ch chan bool
	ch = make(chan bool)

	for i, mapping := range s.Cfg.Link {
		go s.parseLink(mapping, i, ch)
	}

	for i := 0; i < len(s.Cfg.Link); i++ {
		<-ch
	}
}

func (s *Synch) parseMapping(mpngStr string, i int, c chan bool) {
	rawMpng, err := lang.ParseMapping(mpngStr)
	if err != nil {
		panic(err)
	}

	mpng := createMapping(s, rawMpng)
	s.mappings = append(s.mappings, mpng)

	c <- true
}

func (s *Synch) parseMappings() {
	var ch chan bool
	ch = make(chan bool)

	for i, mapping := range s.Cfg.Map {
		go s.parseMapping(mapping, i, ch)
	}

	for i := 0; i < len(s.Cfg.Map); i++ {
		<-ch
	}
}

// selectData selects all records from all tables and filters them to get the relevant records.
func (s *Synch) selectData() {
	for i := range s.Links {
		var lnk *Link = s.Links[i]

		sourceRawActiveRecords := (*lnk.source.db).Select(lnk.source.tbl.name, lnk.sourceWhere)
		targetRawActiveRecords := (*lnk.target.db).Select(lnk.target.tbl.name, lnk.targetWhere)

		// for _, v := range sourceRawActiveRecords {
		// 	fmt.Println(v["film_id"])
		// }

		if !s.initial {
			lnk.sourceOldActiveRecords = lnk.sourceActiveRecords
			lnk.targetOldActiveRecords = lnk.targetActiveRecords
		}

		for _, sourceRecord := range sourceRawActiveRecords {
			sourceRecordPointer := lnk.source.tbl.records.FindRecordPointer(sourceRecord)
			lnk.sourceActiveRecords = append(lnk.sourceActiveRecords, sourceRecordPointer)
			sourceRecordPointer.ActiveIn = append(sourceRecordPointer.ActiveIn, lnk)
		}
		for _, targetRecord := range targetRawActiveRecords {
			targetRecordPointer := lnk.target.tbl.records.FindRecordPointer(targetRecord)
			lnk.targetActiveRecords = append(lnk.targetActiveRecords, targetRecordPointer)
			targetRecordPointer.ActiveIn = append(targetRecordPointer.ActiveIn, lnk)
		}
		// log.Println(lnk.sourceActiveRecords)
		// log.Println(lnk.targetActiveRecords)

	}
}

func (s *Synch) setDatabase(DBMap map[string]*db.Database, dbName string) {
	_, dbExists := DBMap[dbName]
	if dbExists {
		s.dbs[dbName] = DBMap[dbName]
		(*s.dbs[dbName]).Init()
	} else {
		dbErr := &db.DatabaseError{DBName: dbName, ErrMsg: "database hasn't been configured"}
		panic(dbErr)
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
func (s *Synch) Synchronize() {
	for i := range s.Links {
		var lnk *Link = s.Links[i]

		for k := range lnk.pairs {
			var pair *Pair = lnk.pairs[k]
			_, err := pair.Synchronize()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
