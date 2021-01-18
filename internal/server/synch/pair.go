package synch

import (
	"fmt"
	"log"

	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
	"github.com/christoph-karpowicz/db_mediator/internal/server/db"
	"github.com/christoph-karpowicz/db_mediator/internal/util"
)

type pairSynchData struct {
	targetDb        db.Database
	sourceTableName string
	sourceKeyName   string
	sourceKeyValue  interface{}
	targetKeyName   string
	targetExtIDName string
	targetTableName string
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
		link.source.tbl.name,
		link.source.cfg.Key,
		source.Data[link.source.cfg.Key],
		link.target.cfg.Key,
		link.target.matchColumn,
		link.target.tbl.name,
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
	if p.target != nil && util.StringSliceContains(p.Link.synch.GetConfig().Do, cfg.DB_UPDATE) {
		sourceColumnValue := p.source.Data[p.Link.sourceColumn]
		targetColumnValue := p.target.Data[p.Link.targetColumn]

		if areEqual, err := areEqual(sourceColumnValue, targetColumnValue); err != nil {
			log.Println(err)
		} else if !areEqual {
			updateErr := p.doUpdate(sourceColumnValue)
			if updateErr == nil {
				p.logUpdateOrIdleOperation(cfg.OPERATION_UPDATE)
			} else {
				log.Println(updateErr)
			}
		} else {
			if p.Link.synch.GetType() == ONE_OFF && p.Link.synch.IsSimulation() {
				p.logUpdateOrIdleOperation(cfg.OPERATION_IDLE)
			}
		}
	} else if p.target == nil && util.StringSliceContains(p.Link.synch.GetConfig().Do, cfg.DB_INSERT) {
		inDto, insertErr := p.doInsert()
		if insertErr == nil {
			p.logInsertOperation(inDto)
		} else {
			log.Println(insertErr)
		}
	}

	return false, nil
}

func (p Pair) doUpdate(sourceColumnValue interface{}) error {
	upDto := db.UpdateDto{
		p.synchData.sourceTableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		p.Link.targetColumn,
		sourceColumnValue,
	}

	if !p.Link.synch.IsSimulation() {
		update, err := p.synchData.targetDb.Update(upDto)
		if err != nil {
			return err
		}
		log.Println(update)
	}
	return nil
}

func (p Pair) doInsert() (*db.InsertDto, error) {
	inDto := p.prepareInsertValues()
	if !p.Link.synch.IsSimulation() {
		insert, err := p.synchData.targetDb.Insert(*inDto)
		if err != nil {
			return nil, err
		}
		log.Println(insert)
	}
	return inDto, nil
}

func (p *Pair) prepareInsertValues() *db.InsertDto {
	values := make(map[string]interface{})
	for columnName, value := range p.source.Data {
		targetColumn, err := p.findTargetColumnName(columnName)
		if err != nil {
			fmt.Println(err)
		}
		values[targetColumn] = value
	}

	inDto := db.InsertDto{
		p.synchData.sourceTableName,
		p.synchData.targetExtIDName,
		p.synchData.sourceKeyValue,
		values,
	}
	return &inDto
}

func (p *Pair) findTargetColumnName(columnName string) (string, error) {
	for _, mapping := range p.Link.synch.GetRawMappings() {
		if columnName == mapping["sourceColumn"] {
			return mapping["targetColumn"], nil
		}
	}
	return "", &mappingError{errMsg: fmt.Sprintf("Mapping for column \"%s\" not found.", columnName)}
}

func (p *Pair) logUpdateOrIdleOperation(operationType string) {
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

	operation := updateOrIdleOperation{
		Operation:         operationType,
		Timestamp:         util.GetTimestamp(),
		SourceTableName:   p.synchData.sourceTableName,
		SourceKeyName:     p.synchData.sourceKeyName,
		SourceKeyValue:    p.source.Data[p.synchData.sourceKeyName],
		SourceColumnName:  p.Link.sourceColumn,
		SourceColumnValue: sourceColumnValue,
		TargetTableName:   p.synchData.targetTableName,
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

func (p *Pair) logInsertOperation(inDto *db.InsertDto) {
	operation := insertOperation{
		Operation:        cfg.OPERATION_INSERT,
		Timestamp:        util.GetTimestamp(),
		SourceTableName:  p.synchData.sourceTableName,
		SourceKeyName:    p.synchData.sourceKeyName,
		SourceKeyValue:   p.source.Data[p.synchData.sourceKeyName],
		SourceColumnName: p.Link.sourceColumn,
		TargetTableName:  p.synchData.targetTableName,
		InsertedRow:      inDto.Values,
	}

	if !p.Link.synch.IsSimulation() {
		operation.IterationId = p.Link.synch.GetIteration().id
	}

	p.Link.synch.GetIteration().addOperation(&operation)
}
