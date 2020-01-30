package synch

type node struct {
	Name     string `yaml:"name"`
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
	Key      string `yaml:"key"`
}
