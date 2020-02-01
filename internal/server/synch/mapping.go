package synch

type mapping struct {
	source     *node
	target     *node
	matchBy    string
	sourceExID string
	targetExID string
	do         []string
}

func createMapping(str string) *mapping {

	newMapping := mapping{}
	return &newMapping
}

func (m *mapping) parse(str string) {

}
