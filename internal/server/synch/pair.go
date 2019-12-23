package synch

import "errors"

type Pair struct {
	primaryFlow   Flow
	secondaryFlow Flow
}

func CreatePair(rec1 *Record, rec2 *Record, flowSymbol string) (Pair, error) {
	var newPair Pair

	if flowSymbol == "<=" || flowSymbol == "=>" {
		if flowSymbol == "<=" {
			newPair.primaryFlow = Flow{source: rec2, target: rec1}
		} else {
			newPair.primaryFlow = Flow{source: rec1, target: rec2}
		}
		return newPair, nil
	} else if flowSymbol == "*<=>" || flowSymbol == "<=>*" {
		if flowSymbol == "*<=>" {
			newPair.primaryFlow = Flow{source: rec2, target: rec1}
			newPair.secondaryFlow = Flow{source: rec1, target: rec2}
		} else {
			newPair.primaryFlow = Flow{source: rec1, target: rec2}
			newPair.secondaryFlow = Flow{source: rec2, target: rec1}
		}
		return newPair, nil
	}

	return newPair, errors.New("Unknown data flow direction.")
}
