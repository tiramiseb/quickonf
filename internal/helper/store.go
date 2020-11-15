package helper

var store = map[string]string{}

// Store adds a value to the store
func Store(key, value string) {
	store[key] = value
}
