package input

import (
	"fmt"
	"strings"
)

func resolveString(v interface{}, store map[string]interface{}) (str string, err error) {
	switch v.(type) {
	case string:
		str = v.(string)
		for key, val := range store {
			switch val.(type) {
			case string:
				str = strings.ReplaceAll(str, "<"+key+">", val.(string))
			}
		}
	case nil:
		str = ""
	default:
		err = fmt.Errorf("Value \"%v\" is not a string", v)
	}
	return
}
