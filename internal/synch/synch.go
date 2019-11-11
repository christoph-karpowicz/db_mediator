package synch

import (
	"fmt"

	"github.com/christoph-karpowicz/unifier/internal/db"
)

type Synch struct {
	synch     *Data
	database1 *db.Database
	database2 *db.Database
}

func (s *Synch) SetDatabases(DBMap map[string]*db.Database) {
	s.database1 = DBMap[s.synch.Databases[0]]
	s.database2 = DBMap[s.synch.Databases[1]]
	fmt.Println(s)
}
