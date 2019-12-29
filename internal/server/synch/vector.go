package synch

import (
	"fmt"
	"log"
)

type vectorEndpointOptions struct {
	CreateNewRows string `yaml:"createNewRows"`
	UpdateOldRows string `yaml:"updateOldRows"`
}

type vectorEndpoint struct {
	Database  string                `yaml:"database"`
	Table     string                `yaml:"table"`
	Key       string                `yaml:"key"`
	Column    string                `yaml:"column"`
	Condition string                `yaml:"condition"`
	Options   vectorEndpointOptions `yaml:"options"`
}

type vectorSettingsExternalIds struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type vectorSettings struct {
	MatchBy     string                    `yaml:"matchBy"`
	ExternalIds vectorSettingsExternalIds `yaml:"externalIds"`
	Active      bool                      `yaml:"active"`
}

type vector struct {
	Source                 vectorEndpoint `yaml:"source"`
	Target                 vectorEndpoint `yaml:"target"`
	Settings               vectorSettings `yaml:"settings"`
	sourceTable            *table
	targetTable            *table
	sourceOldActiveRecords []*record
	targetOldActiveRecords []*record
	sourceActiveRecords    []*record
	targetActiveRecords    []*record
	pairs                  []pair
}

// For each active record in database1 find a corresponding acitve record in database2.
func (v *vector) createPairs() {
	var sourceRecords []*record
	var targetRecords []*record
	var isBidirectional bool = false

	for i := range sourceRecords {
		source := sourceRecords[i]
		var pairFound bool = false
		for j := range targetRecords {
			target := targetRecords[j]

			if v.Settings.MatchBy == "external_id_columns" {
				var sourceExternalIDColumnName string = v.Settings.ExternalIds.Source
				var targetExternalIDColumnName string = v.Settings.ExternalIds.Target
				sourceExternalID, sourceOk := source.Data[sourceExternalIDColumnName]
				targetExternalID, targetOk := target.Data[targetExternalIDColumnName]

				if !sourceOk || !targetOk {
					continue
				}

				if areEqual, err := areEqual(sourceExternalID, targetExternalID); err != nil {
					log.Println(err)
				} else if areEqual {
					newPair := createPair(v, source, target)
					v.pairs = append(v.pairs, newPair)
					pairFound = true
					source.PairedIn = append(source.PairedIn, v)
					target.PairedIn = append(target.PairedIn, v)
				}
			}
		}
		if !pairFound && isBidirectional {
			newPair := createPair(v, source, nil)
			v.pairs = append(v.pairs, newPair)
		}
	}
	for _, pair := range v.pairs {
		fmt.Printf("rec1: %s\n", pair.source.Data)
		if pair.IsComplete {
			fmt.Printf("rec2: %s\n", pair.target.Data)
		}
		fmt.Println("======")
	}
}
