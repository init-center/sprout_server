package tools

func Contains(sl []interface{}, findItem interface{}) bool {
	for _, item := range sl {
		if item == findItem {
			return true
		}
	}
	return false
}
