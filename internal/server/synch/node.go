package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type node struct {
	cfg         *cfg.NodeConfig
	db          *db.Database
	tbl         *table
	matchColumn string
}

func createNode(cfg *cfg.NodeConfig, db *db.Database, tbl *table) *node {
	newNode := node{
		cfg: cfg,
		db:  db,
		tbl: tbl,
	}
	return &newNode
}

func (n *node) setMatchColumn(col string) {
	n.matchColumn = col
}
