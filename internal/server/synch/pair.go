package synch

import "errors"

type Pair interface {
	GetSource() *Record
	GetTarget() *Record
}

type UnidirectionalPair struct {
	source *Record
	target *Record
}

type BidirectionalPair struct {
	source   *Record
	target   *Record
	Sibling  *BidirectionalPair
	Priority bool
}

func (up UnidirectionalPair) GetSource() *Record {
	return up.source
}

func (up UnidirectionalPair) GetTarget() *Record {
	return up.target
}

func (up BidirectionalPair) GetSource() *Record {
	return up.source
}

func (up BidirectionalPair) GetTarget() *Record {
	return up.target
}

func CreatePair(rec1 *Record, rec2 *Record, flow string) ([2]Pair, error) {
	var pairArray [2]Pair
	if flow == "<=" || flow == "=>" {
		var newPair UnidirectionalPair
		if flow == "<=" {
			newPair = UnidirectionalPair{source: rec2, target: rec1}
		} else {
			newPair = UnidirectionalPair{source: rec1, target: rec2}
		}
		pairArray[0] = newPair
		pairArray[1] = nil
		return pairArray, nil
	} else if flow == "*<=>" || flow == "<=>*" {
		var newPriorityPair BidirectionalPair
		var newNonPriorityPair BidirectionalPair
		if flow == "*<=>" {
			newPriorityPair = BidirectionalPair{source: rec2, target: rec1, Priority: true}
			newNonPriorityPair = BidirectionalPair{source: rec1, target: rec2, Priority: false}
		} else {
			newPriorityPair = BidirectionalPair{source: rec1, target: rec2, Priority: true}
			newNonPriorityPair = BidirectionalPair{source: rec2, target: rec1, Priority: false}
		}
		newPriorityPair.Sibling = &newNonPriorityPair
		newNonPriorityPair.Sibling = &newPriorityPair
		pairArray[0] = newPriorityPair
		pairArray[1] = newNonPriorityPair
		return pairArray, nil
	}
	return pairArray, errors.New("Unknown data flow direction.")
}
