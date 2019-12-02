package synch

type Connection struct {
	Type          string      `json:"type"`
	Columns       []string    `json:"columns"`
	InitialVector interface{} `json:"initial_vector"`
}
