package synch

import (
	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
	"github.com/christoph-karpowicz/db_mediator/internal/server/db"
)

// node holds all the data necessary for
// data exchange.
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
