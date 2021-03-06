package synch

import (
	"log"
	"strings"
	"sync"

	"github.com/christoph-karpowicz/db_mediator/internal/server/cfg"
	"github.com/google/uuid"
)

// Link represents a single link in the config file like:
// [example_node1.example_column1 WHERE ...] TO [example_node2.example_column2 WHERE ...]
type Link struct {
	id           string
	synch        Synchronizer
	Cmd          string
	source       *node
	target       *node
	sourceTable  *table
	targetTable  *table
	sourceColumn string
	targetColumn string
	sourceWhere  string
	targetWhere  string
	sourceExID   string
	targetExID   string
	pairs        []*Pair
}

func createLink(synch Synchronizer, link map[string]string) *Link {

	sourceNode, sourceNodeFound := synch.GetNodes()[link[cfg.PSUBEXP_SOURCE_NODE]]
	if !sourceNodeFound {
		panic("[create link] ERROR: source node not found.")
	}
	targetNode, targetNodeFound := synch.GetNodes()[link[cfg.PSUBEXP_TARGET_NODE]]
	if !targetNodeFound {
		panic("[create link] ERROR: target node not found.")
	}

	newLink := Link{
		id:           uuid.New().String(),
		synch:        synch,
		Cmd:          link["cmd"],
		source:       sourceNode,
		target:       targetNode,
		sourceTable:  sourceNode.tbl,
		targetTable:  targetNode.tbl,
		sourceColumn: link[cfg.PSUBEXP_SOURCE_COLUMN],
		targetColumn: link[cfg.PSUBEXP_TARGET_COLUMN],
		sourceWhere:  link[cfg.PSUBEXP_SOURCE_WHERE],
		targetWhere:  link[cfg.PSUBEXP_TARGET_WHERE],
	}

	if synch.GetConfig().Match.Method == "ids" {
		for _, marg := range synch.GetConfig().Match.Args {
			margSplt := strings.Split(marg, ".")
			margNode := margSplt[0]
			margColumn := margSplt[1]
			if margNode == newLink.source.cfg.Name {
				newLink.sourceExID = margColumn
			} else if margNode == newLink.target.cfg.Name {
				newLink.targetExID = margColumn
			}
		}
	}

	return &newLink
}

func (l Link) GetID() string {
	return l.id
}

func (l *Link) comparePair(src *record, c chan bool) {
	var pairFound bool = false

	for j := range *l.targetTable.activeRecords {
		target := (*l.targetTable.activeRecords)[j]

		if l.synch.GetConfig().Match.Method == "ids" {
			sourceExternalID, sourceOk := src.Data[l.sourceExID]
			targetExternalID, targetOk := target.Data[l.targetExID]

			if !sourceOk || !targetOk {
				continue
			}

			if areEqual, err := areEqual(sourceExternalID, targetExternalID); err != nil {
				log.Println(err)
			} else if areEqual {
				newPair := createPair(l, src, target)
				l.pairs = append(l.pairs, newPair)
				pairFound = true
				src.PairedIn = append(src.PairedIn, l)
				target.PairedIn = append(target.PairedIn, l)
			}
		}
	}

	c <- pairFound
}

// createPairs for each active record in source database finds a corresponding acitve record in target database.
func (l *Link) createPairs(wg *sync.WaitGroup) {
	for i := range *l.sourceTable.activeRecords {
		ch := make(chan bool)
		source := (*l.sourceTable.activeRecords)[i]

		go l.comparePair(source, ch)

		if !<-ch {
			newPair := createPair(l, source, nil)
			l.pairs = append(l.pairs, newPair)
		}
	}
	wg.Done()
}

func (l *Link) reset() {
	l.sourceTable.activeRecords = nil
	l.targetTable.activeRecords = nil
	l.pairs = nil
}
