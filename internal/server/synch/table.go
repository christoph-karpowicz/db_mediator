package synch

type Table struct {
	Names      []string      `json:"names"`
	Keys       []string      `json:"keys"`
	Connection Connection    `json:"connection"`
	Vectors    []interface{} `json:"vectors"`
}
