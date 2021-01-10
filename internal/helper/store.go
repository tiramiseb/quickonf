package helper

var store = map[string]string{}

// Store adds a value to the store. Ignored if key is empty.
func Store(key, value string) {
	if key == "" {
		return
	}
	store[key] = value
}

// AllStore returns the content of the store
func AllStore() map[string]string {
	return store
}
