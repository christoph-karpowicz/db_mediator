package synch

type mapping struct {
	source      *node
	target      *node
	matchMethod string
	sourceExID  string
	targetExID  string
	do          []string
}

func createMapping(nodes map[string]*node, link map[string]string, matchMethod map[string]interface{}, do []string) *mapping {
	newMapping := mapping{
		source:      nodes[link["sourceNode"]],
		target:      nodes[link["targetNode"]],
		matchMethod: matchMethod["matchCmd"].(string),
		sourceExID:  matchMethod["matchArgs"].([]string)[0],
		targetExID:  matchMethod["matchArgs"].([]string)[1],
		do:          do,
	}

	return &newMapping
}
