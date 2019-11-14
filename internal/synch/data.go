package synch

type SynchData struct {
	Name      string   `json:"name"`
	Databases []string `json:"databases"`
	Tables    []Table  `json:"tables"`
}

func (d *SynchData) Validate() {

}
