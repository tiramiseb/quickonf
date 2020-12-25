package helper

import (
	"fmt"
	"strings"
)

// String returns a string with values from store from an interface
func String(v interface{}) (str string, err error) {
	switch v.(type) {
	case string:
		str = v.(string)
		for key, val := range store {
			str = strings.ReplaceAll(str, "<"+key+">", val)
		}
	case nil:
		str = ""
	default:
		err = fmt.Errorf(`Value "%v" is not a string`, v)
	}
	return
}
