package synch

import (
	"log"

	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
	"github.com/christoph-karpowicz/db_mediator/internal/server/db"
)

type dbStore struct {
	nodeCfgs []cfg.NodeConfig
	dbs      map[string]*db.Database
	tables   map[string]*table
	nodes    map[string]*node
}

func (ds *dbStore) Init(DBMap map[string]*db.Database, nodeCfgs []cfg.NodeConfig) {
	ds.dbs = make(map[string]*db.Database)
	ds.tables = make(map[string]*table)
	ds.nodes = make(map[string]*node)

	ds.setNodeCfgs(nodeCfgs)
	ds.setDatabases(DBMap)
	ds.setTables()
	ds.setNodes()
}

func (ds *dbStore) setNodeCfgs(nodeCfgs []cfg.NodeConfig) {
	ds.nodeCfgs = nodeCfgs
}

// setNodes creates node structs and adds them to the relevant synch struct field.
func (ds *dbStore) setNodes() {
	for i := range ds.nodeCfgs {
		var nodeConfig *cfg.NodeConfig = &ds.nodeCfgs[i]

		var tableName string = nodeConfig.Database + "." + nodeConfig.Table
		_, tableFound := ds.tables[tableName]
		if !tableFound {
			log.Fatalln("[create node] ERROR: table " + tableName + " not found.")
		}

		ds.nodes[nodeConfig.Name] = createNode(nodeConfig, ds.dbs[nodeConfig.Database], ds.tables[tableName])
	}
}

// setDatabases opens the chosen database connections.
func (ds *dbStore) setDatabases(DBMap map[string]*db.Database) {
	for j := range ds.nodeCfgs {
		var nodeConfig *cfg.NodeConfig = &ds.nodeCfgs[j]
		ds.setDatabase(DBMap, nodeConfig.Database)
	}
}

func (ds *dbStore) setDatabase(DBMap map[string]*db.Database, dbName string) {
	_, dbExists := DBMap[dbName]
	if dbExists {
		ds.dbs[dbName] = DBMap[dbName]
		(*ds.dbs[dbName]).Init()
	} else {
		dbErr := &db.DatabaseError{DBName: dbName, ErrMsg: "database hasn't been configured"}
		panic(dbErr)
	}
}

// setTables creates table structs based on node yaml data.
func (ds *dbStore) setTables() {
	for j := range ds.nodeCfgs {
		var nodeConfig *cfg.NodeConfig = &ds.nodeCfgs[j]
		ds.setTable(nodeConfig.Table, ds.dbs[nodeConfig.Database])
	}
}

// setTable creates an individual table struct and selects all records from it.
func (ds *dbStore) setTable(tableName string, database *db.Database) {
	var tblID string = (*database).GetConfig().GetName() + "." + tableName
	_, tableSet := ds.tables[tblID]

	if !tableSet {
		tbl := &table{
			id:   tblID,
			db:   database,
			name: tableName,
		}
		ds.tables[tbl.id] = tbl
	}
}
