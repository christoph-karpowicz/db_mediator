package synch

import "github.com/christoph-karpowicz/unifier/internal/server/db"

type table struct {
	id         string
	db         *db.Database
	name       string
	oldRecords *tableRecords
	records    *tableRecords
}
