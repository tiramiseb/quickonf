package helper

func Path(str string) string {
	if len(str) == 0 {
		return Home
	}
	if str[0] == '/' {
		return str
	}
	return Home + "/" + str
}
