package input

import (
	"errors"
)

func MapStringString(input interface{}, store map[string]interface{}) (map[string]string, error) {
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, errors.New("Module configuration must be a map")
	}
	result := make(map[string]string, len(data))
	for k, v := range data {
		str, err := resolveString(v, store)
		if err != nil {
			return nil, err
		}
		key, err := resolveString(k, store)
		if err != nil {
			return nil, err
		}
		result[key] = str
	}
	return result, nil
}
