package synch

type Table struct {
	Names   []string `json:"names"`
	Keys    string   `json:"keys"`
	Vectors []Vector `json:"vectors"`
}
