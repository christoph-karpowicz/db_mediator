package synch

import (
	"regexp"
)

type Vector struct {
	Column1 string
	Dir     string
	Column2 string
}

func (v *Vector) Parse(pair *string) {
	re := regexp.MustCompile(`(\w+)\s([<=>]{2,3})\s(\w+)`)
	res := re.FindStringSubmatch(*pair)
	v.Column1 = res[1]
	v.Dir = res[2]
	v.Column2 = res[3]
}
