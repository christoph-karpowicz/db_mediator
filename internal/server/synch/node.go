package synch

import "github.com/christoph-karpowicz/unifier/internal/server/db"

type nodeData struct {
	Name     string `yaml:"name"`
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Key      string `yaml:"key"`
}

type node struct {
	data *nodeData
	db   *db.Database
	tbl  *table
}

func createNode(data *nodeData, db *db.Database, tbl *table) *node {
	newNode := node{
		data: data,
		db:   db,
		tbl:  tbl,
	}
	return &newNode
}
