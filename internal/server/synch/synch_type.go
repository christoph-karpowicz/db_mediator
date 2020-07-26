package synch

import (
	"errors"
	"strings"
)

type synchType int

const (
	ONE_OFF synchType = iota + 1
	ONGOING
)

// FindSynchType returns a synch type based on a string.
func FindSynchType(sType string) (synchType, error) {
	switch strings.ToLower(sType) {
	case "one-off":
		return ONE_OFF, nil
	case "ongoing":
		return ONGOING, nil
	default:
		return 0, errors.New("Synch type \"" + sType + "\" not found")
	}
}
