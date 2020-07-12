package synch

import (
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// LinkReportData link data for simulation purposes.
type LinkReportData struct {
	Link map[string]string
}

// Link represents a single link in the config file like:
// [example_node1.example_column1 WHERE ...] TO [example_node2.example_column2 WHERE ...]
type Link struct {
	id                     string
	synch                  Synchronizer
	Cmd                    string
	source                 *node
	target                 *node
	sourceColumn           string
	targetColumn           string
	sourceWhere            string
	targetWhere            string
	sourceExID             string
	targetExID             string
	sourceOldActiveRecords records
	sourceActiveRecords    records
	targetOldActiveRecords records
	targetActiveRecords    records
	pairs                  []*Pair
	Rep                    *LinkReportData
}

func createLink(synch Synchronizer, link map[string]string) *Link {

	_, sourceNodeFound := synch.GetNodes()[link["sourceNode"]]
	if !sourceNodeFound {
		panic("[create link] ERROR: source node not found.")
	}
	_, targetNodeFound := synch.GetNodes()[link["targetNode"]]
	if !targetNodeFound {
		panic("[create link] ERROR: target node not found.")
	}

	newLink := Link{
		id:           uuid.New().String(),
		synch:        synch,
		Cmd:          link["cmd"],
		source:       synch.GetNodes()[link["sourceNode"]],
		target:       synch.GetNodes()[link["targetNode"]],
		sourceColumn: link["sourceColumn"],
		targetColumn: link["targetColumn"],
		sourceWhere:  link["sourceWhere"],
		targetWhere:  link["targetWhere"],
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

	newLink.Rep = &LinkReportData{link}

	return &newLink
}

func (l Link) GetID() string {
	return l.id
}

func (l *Link) comparePair(src *record, c chan bool) {
	var pairFound bool = false

	for j := range l.targetActiveRecords {
		target := l.targetActiveRecords[j]

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
	for i := range l.sourceActiveRecords {
		ch := make(chan bool)
		source := l.sourceActiveRecords[i]

		go l.comparePair(source, ch)

		if !<-ch {
			newPair := createPair(l, source, nil)
			l.pairs = append(l.pairs, newPair)
		}
	}
	wg.Done()
	// for _, pair := range l.pairs {
	// 	fmt.Printf("rec1: %s\n", pair.source.Data)
	// 	if pair.target != nil {
	// 		fmt.Printf("rec2: %s\n", pair.target.Data)
	// 	}
	// 	fmt.Println("======")
	// }
}

func (l *Link) reset() {
	l.sourceOldActiveRecords = nil
	l.sourceActiveRecords = nil
	l.targetOldActiveRecords = nil
	l.targetActiveRecords = nil
	l.pairs = nil
	l.Rep = nil
}
