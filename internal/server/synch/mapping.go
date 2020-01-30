package synch

type mapping struct {
	source  string `yaml:"name"`
	Databse string `yaml:"database"`
	Table   string `yaml:"table"`
	Key     string `yaml:"key"`
}
