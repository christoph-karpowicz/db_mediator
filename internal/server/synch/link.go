package synch

import (
	"log"
	"sync"
)

// LinkReportData link data for simulation purposes.
type LinkReportData struct {
	MappingIndex int
	LinkIndex    int
	Link         map[string]string
}

// Link represents a single link in the config file like:
// [example_node1.example_column1 WHERE ...] TO [example_node2.example_column2 WHERE ...]
type Link struct {
	in                     *Instruction
	source                 *node
	target                 *node
	sourceWhere            string
	targetWhere            string
	sourceColumn           string
	targetColumn           string
	matchMethod            string
	sourceExID             string
	targetExID             string
	do                     []string
	sourceOldActiveRecords []*record
	sourceActiveRecords    []*record
	targetOldActiveRecords []*record
	targetActiveRecords    []*record
	pairs                  []*Pair
	Rep                    *LinkReportData
}

func createLink(in *Instruction, link map[string]string, matchMethod map[string]interface{}, do []string, indexes ...int) *Link {

	_, sourceNodeFound := link["sourceNode"]
	if !sourceNodeFound {
		panic("[create link] ERROR: source node not found.")
	}
	_, targetNodeFound := link["targetNode"]
	if !targetNodeFound {
		panic("[create link] ERROR: target node not found.")
	}

	newMapping := Link{
		in:           in,
		source:       in.synch.nodes[link["sourceNode"]],
		target:       in.synch.nodes[link["targetNode"]],
		sourceWhere:  link["sourceWhere"],
		targetWhere:  link["targetWhere"],
		sourceColumn: link["sourceColumn"],
		targetColumn: link["targetColumn"],
		matchMethod:  matchMethod["matchCmd"].(string),
		do:           do,
	}

	if newMapping.matchMethod == "IDS" {
		for _, marg := range matchMethod["parsedMatchArgs"].([]map[string]string) {
			if marg["node"] == newMapping.source.cfg.Name {
				newMapping.sourceExID = marg["extIDColumn"]
			} else {
				newMapping.targetExID = marg["extIDColumn"]
			}
		}
	}

	newMapping.Rep = &LinkReportData{
		MappingIndex: indexes[0],
		LinkIndex:    indexes[1],
		Link:         link,
	}

	return &newMapping
}

func (l *Link) comparePair(src *record, c chan bool) {
	var pairFound bool = false

	for j := range l.targetActiveRecords {
		target := l.targetActiveRecords[j]

		if l.matchMethod == "IDS" {
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
	// 	if pair.IsComplete {
	// 		fmt.Printf("rec2: %s\n", pair.target.Data)
	// 	}
	// 	fmt.Println("======")
	// }
}
