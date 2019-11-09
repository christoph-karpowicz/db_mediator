package db

type Databases struct {
	Databases map[string]Database `json:"databases"`
}

type Database struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
