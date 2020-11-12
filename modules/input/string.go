package input

func String(input interface{}, store map[string]interface{}) (string, error) {
	return resolveString(input, store)
}
