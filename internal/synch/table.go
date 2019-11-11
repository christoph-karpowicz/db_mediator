package synch

type Table struct {
	Names   []string      `json:"names"`
	Keys    []string      `json:"keys"`
	Vectors []interface{} `json:"vectors"`
}
