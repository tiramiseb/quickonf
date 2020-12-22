package helper

import "strings"

var store = map[string]string{}

func init() {
	cmdout, err := Exec(nil, "lsb_release", "--codename", "--short")
	if err != nil {
		panic(err)
	}
	Store("oscodename", strings.TrimSpace(string(cmdout)))
}

// Store adds a value to the store
func Store(key, value string) {
	store[key] = value
}
