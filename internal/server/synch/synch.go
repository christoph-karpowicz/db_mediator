/*
Package synch handles all data sychronization.
*/
package synch

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

// Synch represents an individual synchronzation configration.
// It holds all configuration from an .yaml file, raw and parsed.
type Synch struct {
	cfg        *cfg.SynchConfig
	dbs        map[string]*db.Database
	tables     map[string]*table
	nodes      map[string]*node
	mappings   []*Mapping
	Links      []*Link
	counters   *counters
	stype      string
	running    bool
	initial    bool
	simulation bool
	History    *History
}

// Init prepares the synchronization by fetching all necessary data
// and parsing it.
func (s *Synch) Init(DBMap map[string]*db.Database, stype string) {
	s.History = &History{}
	s.stype = stype
	tStart := time.Now()
	s.dbs = make(map[string]*db.Database)
	s.tables = make(map[string]*table)
	s.nodes = make(map[string]*node)

	if s.counters == nil {
		s.counters = newCounters()
		s.setDatabases(DBMap)
		s.setTables()
		s.setNodes()
		s.parseCfgLinks()
		s.parseCfgMappings()
		s.parseCfgMatcher()
	}

	fmt.Println("Synch init finished in: ", time.Since(tStart).String())
}

// GetConfig returns the synch config struct.
func (s *Synch) GetConfig() *cfg.SynchConfig {
	return s.cfg
}

// GetHistory returns the synch's history.
func (s *Synch) GetHistory() *History {
	return s.History
}

// GetNodes returns all nodes between which
// synchronization takes place.
func (s *Synch) GetNodes() map[string]*node {
	return s.nodes
}

// GetType returns the type of the synch.
func (s *Synch) GetType() string {
	return s.stype
}

func (s *Synch) IsInitial() bool {
	return s.initial
}

// SetInitial sets the initial struct field indicating whether
// it's the first run of the synch.
func (s *Synch) SetInitial(ini bool) {
	s.initial = ini
}

func (s *Synch) IsRunning() bool {
	return s.running
}

func (s *Synch) IsSimulation() bool {
	return s.simulation
}

func (s *Synch) SetSimulation(sim bool) {
	s.simulation = sim
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

func (s *Synch) parseCfgLinks() {
	var ch chan bool
	ch = make(chan bool)

	for i, mapping := range s.cfg.Link {
		go s.parseLink(mapping, i, ch)
	}

	for i := 0; i < len(s.cfg.Link); i++ {
		<-ch
	}
}

func (s *Synch) parseLink(mpngStr string, i int, c chan bool) {
	rawLink, err := cfg.ParseLink(mpngStr)
	if err != nil {
		panic(err)
	}

	in := createLink(s, rawLink)
	s.Links = append(s.Links, in)

	c <- true
}

func (s *Synch) parseCfgMappings() {
	var ch chan bool
	ch = make(chan bool)

	for i, mapping := range s.cfg.Map {
		go s.parseMapping(mapping, i, ch)
	}

	for i := 0; i < len(s.cfg.Map); i++ {
		<-ch
	}
}

func (s *Synch) parseMapping(mpngStr string, i int, c chan bool) {
	rawMpng, err := cfg.ParseMapping(mpngStr)
	if err != nil {
		panic(err)
	}

	mpng := createMapping(s, rawMpng)
	s.mappings = append(s.mappings, mpng)

	c <- true
}

func (s *Synch) parseCfgMatcher() {
	matcherMethod := s.GetConfig().Match.Method

	switch matcherMethod {
	case "ids":
		parsedMatcher, err := cfg.ParseIdsMatcherMethod(s.GetConfig().Match.Args)
		if err != nil {
			panic(err)
		}

		for _, arg := range parsedMatcher {
			node, found := s.nodes[arg[0]]
			if !found {
				panic(errors.New("node name not found"))
			}

			node.setMatchColumn(arg[1])
		}
	default:
		panic(errors.New("unknown match method"))
	}
}

// selectData selects all records from all tables and filters them to get the relevant records.
func (s *Synch) selectData() {
	for i := range s.Links {
		var lnk *Link = s.Links[i]

		sourceRawActiveRecords := (*lnk.source.db).Select(lnk.source.tbl.name, lnk.sourceWhere)
		targetRawActiveRecords := (*lnk.target.db).Select(lnk.target.tbl.name, lnk.targetWhere)

		// if !s.initial {
		// 	lnk.sourceOldActiveRecords = lnk.sourceActiveRecords
		// 	lnk.targetOldActiveRecords = lnk.targetActiveRecords
		// }

		lnk.sourceTable.setActiveRecords(sourceRawActiveRecords)
		lnk.targetTable.setActiveRecords(targetRawActiveRecords)

		lnk.sourceTable.activeRecords.setActiveIn(lnk)
		lnk.targetTable.activeRecords.setActiveIn(lnk)
	}

	s.counters.selects++
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
	for j := range s.cfg.Nodes {
		var nodeConfig *cfg.NodeConfig = &s.cfg.Nodes[j]
		s.setDatabase(DBMap, nodeConfig.Database)
	}
}

// setNodes creates node structs and adds them to the relevant synch struct field.
func (s *Synch) setNodes() {
	for i := range s.cfg.Nodes {
		var nodeConfig *cfg.NodeConfig = &s.cfg.Nodes[i]

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
	_, tableSet := s.tables[tblID]

	if !tableSet {
		tbl := &table{
			id:   tblID,
			db:   database,
			name: tableName,
		}
		s.tables[tbl.id] = tbl
	}
}

// setTables creates table structs based on node yaml data.
func (s *Synch) setTables() {
	for j := range s.cfg.Nodes {
		var nodeConfig *cfg.NodeConfig = &s.cfg.Nodes[j]
		s.setTable(nodeConfig.Table, s.dbs[nodeConfig.Database])
	}
}

// Run executes a single run of the synchronization.
func (s *Synch) Run() {
	s.running = true

	s.selectData()
	s.pairData()
	s.synchronize()
	s.flush()
}

// Stop stops the synch.
func (s *Synch) Stop() {
	s.running = false
}

// synchronize loops over all pairs in all mappings and invokes their synchronize function.
func (s *Synch) synchronize() {
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

func (s *Synch) flush() {
	for i := range s.Links {
		s.Links[i].flush()
	}
}

// Reset clears data preparing the Synch for the next run.
func (s *Synch) Reset() {
	s.stype = ""
	s.SetInitial(false)
	for _, lnk := range s.Links {
		lnk.reset()
	}
	s.counters.reset()
}
