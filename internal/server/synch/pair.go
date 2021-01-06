package synch

import (
	"log"
	"time"

	"github.com/christoph-karpowicz/unifier/internal/server/cfg"
	"github.com/christoph-karpowicz/unifier/internal/server/db"
	arrayUtil "github.com/christoph-karpowicz/unifier/internal/util/array"
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
// Updates if this pair is complete (has both the source and the target)
// and inserts if a target record has to be created.
func (p Pair) Synchronize() (bool, error) {
	if p.target != nil && arrayUtil.Contains(p.Link.synch.GetConfig().Do, cfg.DB_UPDATE) {
		sourceColumnValue := p.source.Data[p.Link.sourceColumn]
		targetColumnValue := p.target.Data[p.Link.targetColumn]

		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
			log.Println(err)
		} else if !areEqual {
			var updateErr error
			if !p.Link.synch.IsSimulation() {
				// updateErr = p.doUpdate(sourceColumnValue)
			}

			if updateErr == nil {
				p.logAction(cfg.OPERATION_UPDATE)
			} else {
				log.Println(updateErr)
			}
		} else {
			if p.Link.synch.GetType() == ONE_OFF {
				p.logAction(cfg.OPERATION_IDLE)
			}
		}
	} else if p.target == nil && arrayUtil.Contains(p.Link.synch.GetConfig().Do, cfg.DB_INSERT) {
		var insertErr error
		if !p.Link.synch.IsSimulation() {
			// p.doInsert()
		}

		if insertErr == nil {
			p.logAction(cfg.OPERATION_INSERT)
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
func (p *Pair) logAction(operationType string) {
	var targetKeyName string
	var targetKeyValue interface{}
	var targetColumnValue interface{}

	var sourceColumnValue interface{} = p.source.Data[p.Link.sourceColumn].(interface{})
	if p.target != nil {
		targetKeyValue = p.target.Data[p.synchData.targetKeyName]
		targetKeyName = p.synchData.targetKeyName
		targetColumnValue = p.target.Data[p.Link.targetColumn].(interface{})
	} else {
		targetKeyName = ""
		targetKeyValue = nil
		targetColumnValue = nil
	}

	dateLayout := "Mon, 02 Jan 2006 15:04:05 MST"
	date := time.Now()

	operation := operation{
		Operation:         operationType,
		Timestamp:         date.Format(dateLayout),
		SourceKeyName:     p.synchData.sourceKeyName,
		SourceKeyValue:    p.source.Data[p.synchData.sourceKeyName],
		SourceColumnName:  p.Link.sourceColumn,
		SourceColumnValue: sourceColumnValue,
		TargetKeyName:     targetKeyName,
		TargetKeyValue:    targetKeyValue,
		TargetColumnName:  p.Link.targetColumn,
		TargetColumnValue: targetColumnValue,
	}

	if !p.Link.synch.IsSimulation() {
		operation.IterationId = p.Link.synch.GetIteration().id
	}

	p.Link.synch.GetIteration().addOperation(&operation)
}
