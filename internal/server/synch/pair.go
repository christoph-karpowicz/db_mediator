package synch

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/christoph-karpowicz/unifier/internal/server/db"
	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

type pairSynchData struct {
	targetDb        db.Database
	tableName       string
	sourceKeyName   string
	sourceKeyValue  interface{}
	targetKeyName   string
	targetExtIDName string
}

// Pair represents a connection between two records, that are going
// to be synchronized.
// Can be complete or incomplete, where incomplete means that there's
// only a source record and the target record will have to be created
// if the synchronization is to be carried out.
// When a pair is incomplete, a target record will be created only if the
// parent mapping is configured to DO INSERTs.
type Pair struct {
	Link      *Link
	source    *record
	target    *record
	synchData *pairSynchData
}

func createPair(link *Link, source *record, target *record) *Pair {
	var synchData pairSynchData = pairSynchData{
		*link.target.db,
		link.target.tbl.name,
		link.source.cfg.Key,
		source.Data[link.source.cfg.Key],
		link.target.cfg.Key,
		link.synch.GetConfig().MatchBy.Args[1],
	}

	var newPair Pair = Pair{
		link,
		source,
		target,
		&synchData,
	}

	return &newPair
}

// Synchronize carries out the synchronization of the two records.
func (p Pair) Synchronize() (bool, error) {
	if p.target != nil && arrUtil.Contains(p.Link.synch.GetConfig().Do, "UPDATE") {
		// Updates
		// If this pair is complete.
		// log.Println(p.source)
		// log.Println(p.target)

		sourceColumnValue := p.source.Data[p.Link.sourceColumn]
		targetColumnValue := p.target.Data[p.Link.targetColumn]

		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
			log.Println(err)
		} else if !areEqual {
			// fmt.Println(sourceColumnValue)
			// fmt.Println(targetColumnValue)

			if !p.Link.synch.Simulation {
				p.doUpdate(sourceColumnValue)
			}

			_, err := p.Link.synch.Rep.AddAction(p, "update")
			if err != nil {
				panic(err)
			}
			// fmt.Println(p.Link.synch.Simulation)
		} else {
			_, err := p.Link.synch.Rep.AddAction(p, "idle")
			if err != nil {
				panic(err)
			}
		}
	} else if p.target == nil && arrUtil.Contains(p.Link.synch.GetConfig().Do, "INSERT") {
		// Inserts
		// If a target record has to be created.
		// log.Println(p.source)
		// log.Println(p.target)

		if !p.Link.synch.Simulation {
			p.doInsert()
		}

		_, err := p.Link.synch.Rep.AddAction(p, "insert")
		if err != nil {
			panic(err)
		}
	}

	return false, nil
}

func (p Pair) doUpdate(sourceColumnValue interface{}) {
	upDto := db.UpdateDto{
		p.synchData.tableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		p.Link.targetColumn,
		sourceColumnValue,
	}

	update, err := p.synchData.targetDb.Update(upDto)
	if err != nil {
		log.Println(err)
	}
	log.Println(update)
	// log.Println(sourceColumnValue)
	// log.Println(targetColumnValue)
}

func (p Pair) doInsert() {
	inDto := db.InsertDto{
		p.synchData.tableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		p.source.Data,
	}

	insert, err := p.synchData.targetDb.Insert(inDto)
	if err != nil {
		log.Println(err)
	}
	log.Println(insert)
}

// ReportJSON creates a JSON representation of an action.
func (p Pair) ReportJSON(actionType string) ([]byte, error) {
	var sourceColumnData interface{} = p.source.Data[p.Link.sourceColumn].(interface{})
	if reflect.TypeOf(sourceColumnData).Name() == "string" && len(sourceColumnData.(string)) > 25 {
		sourceColumnData = sourceColumnData.(string)[:22] + "..."
	}

	var targetKeyName string
	var targetKeyValue interface{}
	var targetColumnData interface{}

	if p.target != nil {
		targetKeyValue = p.target.Data[p.synchData.targetKeyName]
		targetKeyName = p.synchData.targetKeyName

		targetColumnData = p.target.Data[p.Link.targetColumn].(interface{})
		if reflect.TypeOf(targetColumnData).Name() == "string" && len(targetColumnData.(string)) > 25 {
			targetColumnData = targetColumnData.(string)[:22] + "..."
		}
	} else {
		targetKeyName = ""
		targetKeyValue = nil
		targetColumnData = nil
	}

	actionStruct := struct {
		SourceNodeKey    string      `json:"sourceNodeKey"`
		SourceData       interface{} `json:"sourceData"`
		SourceColumn     string      `json:"sourceColumn"`
		SourceColumnData interface{} `json:"sourceColumnData"`
		TargetKeyName    string      `json:"targetKeyName"`
		TargetKeyValue   interface{} `json:"targetKeyValue"`
		TargetColumn     string      `json:"targetColumn"`
		TargetColumnData interface{} `json:"targetColumnData"`
		ActionType       string      `json:"actionType"`
	}{
		SourceNodeKey:    p.synchData.sourceKeyName,
		SourceData:       p.source.Data[p.synchData.sourceKeyName],
		SourceColumn:     p.Link.sourceColumn,
		SourceColumnData: sourceColumnData,
		TargetKeyName:    targetKeyName,
		TargetKeyValue:   targetKeyValue,
		TargetColumn:     p.Link.targetColumn,
		TargetColumnData: targetColumnData,
		ActionType:       actionType,
	}

	return json.Marshal(&actionStruct)
}
