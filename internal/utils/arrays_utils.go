package utils

func AllEqual(array []string) bool {
	if len(array) == 0 {
		return false
	}

	base := array[0]

	for _, v := range array[1:] {
		if v != base {
			return false
		}
	}

	return true
}
