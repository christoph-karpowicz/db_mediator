package cfg

type Match struct {
	Method string   `yaml:"method"`
	Args   []string `yaml:"args"`
}
