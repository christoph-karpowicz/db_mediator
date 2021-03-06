package synch

import "github.com/christoph-karpowicz/db_mediator/internal/server/db"

type table struct {
	id            string
	db            *db.Database
	name          string
	activeRecords *records
}

func (t *table) setActiveRecords(records []map[string]interface{}) {
	t.activeRecords = mapToRecords(records)
}
