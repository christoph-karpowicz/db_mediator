package synch

import (
	"encoding/json"
	"log"
	"reflect"

	arrUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
)

// Pair represents a connection between two records, that are going
// to be synchronized.
// Can be complete or incomplete, where incomplete means that there's
// only a source record and the target record will have to be created
// if the synchronization is to be carried out.
// When a pair is incomplete, a target record will be created only if the
// parent mapping is configured to DO INSERTs.
type Pair struct {
	Mapping *Mapping
	source  *record
	target  *record
}

func createPair(mpng *Mapping, source *record, target *record) *Pair {
	var newPair Pair = Pair{
		Mapping: mpng,
		source:  source,
		target:  target,
	}

	return &newPair
}

func (p Pair) getSourceNodeKey() string {
	return p.Mapping.source.cfg.Key
}

func (p Pair) getTargetNodeKey() string {
	return p.Mapping.target.cfg.Key
}

// Synchronize carries out the synchronization of the two records.
func (p Pair) Synchronize() (bool, error) {
	// db1 := p.Mapping.source.db
	// db2 := p.Mapping.target.db

	if p.target != nil && arrUtil.Contains(p.Mapping.do, "UPDATE") {
		// Updates
		// If this pair is complete.
		// log.Println(p.source)
		// log.Println(p.target)

		sourceColumnValue := p.source.Data[p.Mapping.sourceColumn]
		targetColumnValue := p.target.Data[p.Mapping.targetColumn]

		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
			log.Println(err)
		} else if !areEqual {
			if !p.Mapping.synch.Simulation {
				update, err := (*p.Mapping.target.db).Update(p.Mapping.target.tbl.name, p.getTargetNodeKey(), p.target.Data[p.getTargetNodeKey()], p.Mapping.targetColumn, sourceColumnValue)
				if err != nil {
					log.Println(err)
				}
				log.Println(update)
				// log.Println(sourceColumnValue)
				// log.Println(targetColumnValue)
			}

			_, err := p.Mapping.synch.Rep.AddAction(p, "update")
			if err != nil {
				panic(err)
			}
			// fmt.Println(p.Mapping.synch.Simulation)
		} else {
			_, err := p.Mapping.synch.Rep.AddAction(p, "idle")
			if err != nil {
				panic(err)
			}
		}
	} else if p.target == nil && arrUtil.Contains(p.Mapping.do, "INSERT") {
		// Inserts
		// If a target record has to be created.
		if !p.Mapping.synch.Simulation {
		}

		_, err := p.Mapping.synch.Rep.AddAction(p, "insert")
		if err != nil {
			panic(err)
		}
		// fmt.Println(p.Mapping.synch.Simulation)
	}

	return false, nil
}

// ReportJSON creates a JSON representation of an action.
func (p Pair) ReportJSON(actionType string) ([]byte, error) {
	var sourceColumnData interface{} = p.source.Data[p.Mapping.sourceColumn].(interface{})
	if reflect.TypeOf(sourceColumnData).Name() == "string" && len(sourceColumnData.(string)) > 25 {
		sourceColumnData = sourceColumnData.(string)[:22] + "..."
	}

	var targetKeyName string
	var targetKeyValue interface{}
	var targetColumnData interface{}

	if p.target != nil {
		targetKeyValue = p.target.Data[p.getTargetNodeKey()]
		targetKeyName = p.getTargetNodeKey()

		targetColumnData = p.target.Data[p.Mapping.targetColumn].(interface{})
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
		SourceNodeKey:    p.getSourceNodeKey(),
		SourceData:       p.source.Data[p.getSourceNodeKey()],
		SourceColumn:     p.Mapping.sourceColumn,
		SourceColumnData: sourceColumnData,
		TargetKeyName:    targetKeyName,
		TargetKeyValue:   targetKeyValue,
		TargetColumn:     p.Mapping.targetColumn,
		TargetColumnData: targetColumnData,
		ActionType:       actionType,
	}

	return json.Marshal(&actionStruct)
}
