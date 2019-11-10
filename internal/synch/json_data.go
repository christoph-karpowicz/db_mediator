package synch

type Data struct {
	Name      string   `json:"name"`
	Databases []string `json:"databases"`
	Tables    []Table  `json:"tables"`
}
