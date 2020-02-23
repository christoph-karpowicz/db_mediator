package synch

import "github.com/christoph-karpowicz/unifier/internal/server/db"

type nodeConfig struct {
	Name     string `yaml:"name"`
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Key      string `yaml:"key"`
}

type node struct {
	cfg *nodeConfig
	db  *db.Database
	tbl *table
}

func createNode(cfg *nodeConfig, db *db.Database, tbl *table) *node {
	newNode := node{
		cfg: cfg,
		db:  db,
		tbl: tbl,
	}
	return &newNode
}
