package synch

type matchBy struct {
	Method string   `yaml:"method"`
	Args   []string `yaml:"args"`
}

func createMatchBy() *matchBy {
	newMatchBy := matchBy{}
	return &newMatchBy
}
