package synch

import (
	"database/sql"
	"fmt"
	"log"

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
	if DBMap[s.synch.Databases.Db1.Name] == nil || DBMap[s.synch.Databases.Db2.Name] == nil {
		panic(s.synch.Name + " database config is invalid.")
	}
	s.database1 = DBMap[s.synch.Databases.Db1.Name]
	s.database2 = DBMap[s.synch.Databases.Db2.Name]
}

func (s *Synch) SelectData() {
	var db1Data *db.DatabaseData = (*s.database1).GetData()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db1Data.Host, db1Data.Port, db1Data.User, db1Data.Password, db1Data.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	rows, err := db.Query(`SELECT film_id, title FROM film WHERE title ILIKE 'Des%'`)
	for rows.Next() {
		var id string
		var title string

		if err := rows.Scan(&id, &title); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %s, title: %s\n", id, title)
	}
}
