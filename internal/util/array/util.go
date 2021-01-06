package arrayUtil

func Contains(array []string, val string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}
