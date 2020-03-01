package synch

import (
	"log"
	"sync"
)

// MappingReportData mapping data for simulation purposes.
type MappingReportData struct {
	MappingIndex int
	LinkIndex    int
	Link         map[string]string
}

// Mapping represents a single link in the config file like:
// [example_node1.example_column1 WHERE ...] TO [example_node2.example_column2 WHERE ...]
// It contains all data from the surrounding MAP command.
// There is no struct representing the whole MAP command.
type Mapping struct {
	synch                  *Synch
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
	Rep                    *MappingReportData
}

func createMapping(synch *Synch, link map[string]string, matchMethod map[string]interface{}, do []string, indexes ...int) *Mapping {

	_, sourceNodeFound := link["sourceNode"]
	if !sourceNodeFound {
		log.Fatalln("[create mapping] ERROR: source node not found.")
	}
	_, targetNodeFound := link["targetNode"]
	if !targetNodeFound {
		log.Fatalln("[create mapping] ERROR: target node not found.")
	}

	newMapping := Mapping{
		synch:        synch,
		source:       synch.nodes[link["sourceNode"]],
		target:       synch.nodes[link["targetNode"]],
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

	newMapping.Rep = &MappingReportData{
		MappingIndex: indexes[0],
		LinkIndex:    indexes[1],
		Link:         link,
	}

	return &newMapping
}

func (m *Mapping) comparePair(src *record, c chan bool) {
	var pairFound bool = false

	for j := range m.targetActiveRecords {
		target := m.targetActiveRecords[j]

		if m.matchMethod == "IDS" {
			sourceExternalID, sourceOk := src.Data[m.sourceExID]
			targetExternalID, targetOk := target.Data[m.targetExID]

			if !sourceOk || !targetOk {
				continue
			}

			if areEqual, err := areEqual(sourceExternalID, targetExternalID); err != nil {
				log.Println(err)
			} else if areEqual {
				newPair := createPair(m, src, target)
				m.pairs = append(m.pairs, newPair)
				pairFound = true
				src.PairedIn = append(src.PairedIn, m)
				target.PairedIn = append(target.PairedIn, m)
			}
		}
	}

	c <- pairFound
}

// createPairs for each active record in source database finds a corresponding acitve record in target database.
func (m *Mapping) createPairs(wg *sync.WaitGroup) {
	for i := range m.sourceActiveRecords {
		ch := make(chan bool)
		source := m.sourceActiveRecords[i]

		go m.comparePair(source, ch)

		if !<-ch {
			newPair := createPair(m, source, nil)
			m.pairs = append(m.pairs, newPair)
		}
	}
	wg.Done()
	// for _, pair := range m.pairs {
	// 	fmt.Printf("rec1: %s\n", pair.source.Data)
	// 	if pair.IsComplete {
	// 		fmt.Printf("rec2: %s\n", pair.target.Data)
	// 	}
	// 	fmt.Println("======")
	// }
}
