package synch

import (
	"log"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
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
		link.target.matchColumn,
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
	if p.target != nil && arrUtil.Contains(p.Link.synch.GetConfig().Do, cfg.DB_UPDATE) {
		// Updates if this pair is complete.
		sourceColumnValue := p.source.Data[p.Link.sourceColumn]
		targetColumnValue := p.target.Data[p.Link.targetColumn]

		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
			log.Println(err)
		} else if !areEqual {
			var updateErr error
			if !p.Link.synch.IsSimulation() {
				updateErr = p.doUpdate(sourceColumnValue)
			}

			if updateErr == nil {
				p.logAction(cfg.UPDATE_ACTION)
			} else {
				log.Println(updateErr)
			}
		} else {
			if p.Link.synch.GetType() == ONE_OFF {
				p.logAction(cfg.IDLE_ACTION)
			}
		}
	} else if p.target == nil && arrUtil.Contains(p.Link.synch.GetConfig().Do, cfg.DB_INSERT) {
		// Inserts
		// If a target record has to be created.
		var insertErr error
		if !p.Link.synch.IsSimulation() {
			p.doInsert()
		}

		if insertErr == nil {
			p.logAction(cfg.INSERT_ACTION)
		} else {
			log.Println(insertErr)
		}
	}

	return false, nil
}

func (p Pair) doUpdate(sourceColumnValue interface{}) error {
	upDto := db.UpdateDto{
		p.synchData.tableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		p.Link.targetColumn,
		sourceColumnValue,
	}

	update, err := p.synchData.targetDb.Update(upDto)
	if err != nil {
		return err
	}
	log.Println(update)
	return nil
}

func (p Pair) doInsert() error {
	inDto := db.InsertDto{
		p.synchData.tableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		p.source.Data,
	}

	insert, err := p.synchData.targetDb.Insert(inDto)
	if err != nil {
		return err
	}
	log.Println(insert)
	return nil
}

// logAction adds an action to synch history.
func (p *Pair) logAction(actType string) {
	var sourceColumnData interface{} = p.source.Data[p.Link.sourceColumn].(interface{})
	// if reflect.TypeOf(sourceColumnData).Name() == "string" && len(sourceColumnData.(string)) > 25 {
	// 	sourceColumnData = sourceColumnData.(string)[:22] + "..."
	// }

	var targetKeyName string
	var targetKeyValue interface{}
	var targetColumnData interface{}

	if p.target != nil {
		targetKeyValue = p.target.Data[p.synchData.targetKeyName]
		targetKeyName = p.synchData.targetKeyName

		targetColumnData = p.target.Data[p.Link.targetColumn].(interface{})
		// if reflect.TypeOf(targetColumnData).Name() == "string" && len(targetColumnData.(string)) > 25 {
		// 	targetColumnData = targetColumnData.(string)[:22] + "..."
		// }
	} else {
		targetKeyName = ""
		targetKeyValue = nil
		targetColumnData = nil
	}

	act := action{
		linkId:           p.Link.GetID(),
		ActType:          actType,
		SourceNodeKey:    p.synchData.sourceKeyName,
		SourceData:       p.source.Data[p.synchData.sourceKeyName],
		SourceColumn:     p.Link.sourceColumn,
		SourceColumnData: sourceColumnData,
		TargetKeyName:    targetKeyName,
		TargetKeyValue:   targetKeyValue,
		TargetColumn:     p.Link.targetColumn,
		TargetColumnData: targetColumnData,
	}

	p.Link.synch.GetIteration().addAction(&act)
	// p.Link.synch.GetHistory().addAction(&act)
}
