package helper

func IsInStrings(str string, strs []string) bool {
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}
