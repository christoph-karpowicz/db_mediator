package cfg

type MatchBy struct {
	Method string   `yaml:"method"`
	Args   []string `yaml:"args"`
}
