package synch

import (
	"log"
)

type mapping struct {
	synch                  *synch
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
	pairs                  []*pair
}

func createMapping(synch *synch, link map[string]string, matchMethod map[string]interface{}, do []string) *mapping {
	_, sourceNodeFound := link["sourceNode"]
	if !sourceNodeFound {
		log.Fatalln("[create mapping] ERROR: source node not found.")
	}
	_, targetNodeFound := link["targetNode"]
	if !targetNodeFound {
		log.Fatalln("[create mapping] ERROR: target node not found.")
	}

	newMapping := mapping{
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
			if marg["node"] == newMapping.source.data.Name {
				newMapping.sourceExID = marg["extIDColumn"]
			} else {
				newMapping.targetExID = marg["extIDColumn"]
			}
		}
	}

	return &newMapping
}

// createPairs for each active record in source database finds a corresponding acitve record in target database.
func (m *mapping) createPairs() {
	for i := range m.source.tbl.records.records {
		source := &m.source.tbl.records.records[i]
		var pairFound bool = false

		for j := range m.target.tbl.records.records {
			target := &m.target.tbl.records.records[j]

			if m.matchMethod == "IDS" {
				sourceExternalID, sourceOk := source.Data[m.sourceExID]
				targetExternalID, targetOk := target.Data[m.targetExID]

				if !sourceOk || !targetOk {
					continue
				}

				if areEqual, err := areEqual(sourceExternalID, targetExternalID); err != nil {
					log.Println(err)
				} else if areEqual {
					newPair := createPair(m, source, target)
					m.pairs = append(m.pairs, newPair)
					pairFound = true
					source.PairedIn = append(source.PairedIn, m)
					target.PairedIn = append(target.PairedIn, m)
				}
			}
		}
		if !pairFound {
			newPair := createPair(m, source, nil)
			m.pairs = append(m.pairs, newPair)
		}
	}
	// for _, pair := range m.pairs {
	// 	fmt.Printf("rec1: %s\n", pair.source.Data)
	// 	if pair.IsComplete {
	// 		fmt.Printf("rec2: %s\n", pair.target.Data)
	// 	}
	// 	fmt.Println("======")
	// }
}
