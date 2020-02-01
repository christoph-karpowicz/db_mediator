package lang

type RawMapping struct {
	Links   []string
	MatchBy string
	Do      string
}

func Parse(str string) RawMapping {
	rawMpng := RawMapping{}
	return rawMpng
}
