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
	Link   *Link
	source *record
	target *record
}

func createPair(mpng *Link, source *record, target *record) *Pair {
	var newPair Pair = Pair{
		Link:   mpng,
		source: source,
		target: target,
	}

	return &newPair
}

func (p Pair) getSourceNodeKey() string {
	return p.Link.source.cfg.Key
}

func (p Pair) getTargetNodeKey() string {
	return p.Link.target.cfg.Key
}

// Synchronize carries out the synchronization of the two records.
func (p Pair) Synchronize() (bool, error) {

	if p.target != nil && arrUtil.Contains(p.Link.synch.Cfg.Do, "UPDATE") {
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

			// if !p.Link.In.synch.Simulation {
			// 	update, err := (*p.Link.target.db).Update(p.Link.target.tbl.name, p.getTargetNodeKey(), p.target.Data[p.getTargetNodeKey()], p.Link.targetColumn, sourceColumnValue)
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	log.Println(update)
			// 	// log.Println(sourceColumnValue)
			// 	// log.Println(targetColumnValue)
			// }

			// _, err := p.Link.In.synch.Rep.AddAction(p, "update")
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println(p.Link.In.synch.Simulation)
		} else {
			// _, err := p.Link.In.synch.Rep.AddAction(p, "idle")
			// if err != nil {
			// 	panic(err)
			// }
		}
	} else if p.target == nil && arrUtil.Contains(p.Link.synch.Cfg.Do, "INSERT") {
		// Inserts
		// If a target record has to be created.
		// log.Println(p.source)
		// log.Println(p.target)

		// if !p.Link.In.synch.Simulation {
		// }

		// _, err := p.Link.In.synch.Rep.AddAction(p, "insert")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(p.Link.In.synch.Simulation)
	}

	return false, nil
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
		targetKeyValue = p.target.Data[p.getTargetNodeKey()]
		targetKeyName = p.getTargetNodeKey()

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
		SourceNodeKey:    p.getSourceNodeKey(),
		SourceData:       p.source.Data[p.getSourceNodeKey()],
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
