package synch

import "log"

type mapping struct {
	source                 *node
	target                 *node
	sourceWhere            string
	targetWhere            string
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

func createMapping(nodes map[string]*node, link map[string]string, matchMethod map[string]interface{}, do []string) *mapping {
	_, sourceNodeFound := link["sourceNode"]
	if !sourceNodeFound {
		log.Fatalln("[create mapping] ERROR: source node not found.")
	}
	_, targetNodeFound := link["targetNode"]
	if !targetNodeFound {
		log.Fatalln("[create mapping] ERROR: target node not found.")
	}

	newMapping := mapping{
		source:      nodes[link["sourceNode"]],
		target:      nodes[link["targetNode"]],
		sourceWhere: link["sourceWhere"],
		targetWhere: link["targetWhere"],
		matchMethod: matchMethod["matchCmd"].(string),
		sourceExID:  matchMethod["matchArgs"].([]string)[0],
		targetExID:  matchMethod["matchArgs"].([]string)[1],
		do:          do,
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
				var sourceExternalIDColumnName string = m.sourceExID
				var targetExternalIDColumnName string = m.targetExID
				sourceExternalID, sourceOk := source.Data[sourceExternalIDColumnName]
				targetExternalID, targetOk := target.Data[targetExternalIDColumnName]

				if !sourceOk || !targetOk {
					continue
				}

				if areEqual, err := areEqual(sourceExternalID, targetExternalID); err != nil {
					log.Println(err)
				} else if areEqual {
					newPair := createPair(m, source, target)
					m.pairs = append(m.pairs, &newPair)
					pairFound = true
					source.PairedIn = append(source.PairedIn, m)
					target.PairedIn = append(target.PairedIn, m)
				}
			}
		}
		if !pairFound {
			newPair := createPair(m, source, nil)
			m.pairs = append(m.pairs, &newPair)
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
