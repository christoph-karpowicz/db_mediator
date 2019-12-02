package synch

import (
	"github.com/christoph-karpowicz/unifier/internal/server/db"
)

type Synch struct {
	synch     *SynchData
	database1 *db.Database
	database2 *db.Database
}

func (s *Synch) GetData() *SynchData {
	return s.synch
}

func (s *Synch) SetDatabases(DBMap map[string]*db.Database) {
	if len(s.synch.Databases) < 2 || DBMap[s.synch.Databases[0]] == nil || DBMap[s.synch.Databases[1]] == nil {
		panic(s.synch.Name + " database config is invalid.")
	}
	s.database1 = DBMap[s.synch.Databases[0]]
	s.database2 = DBMap[s.synch.Databases[1]]
}
