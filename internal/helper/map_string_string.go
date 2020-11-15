package helper

import (
	"errors"
)

// MapStringString transforms an interface to a map from strings to strings
func MapStringString(input interface{}) (map[string]string, error) {
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, errors.New("Module configuration must be a map")
	}
	result := make(map[string]string, len(data))
	for k, v := range data {
		str, err := String(v)
		if err != nil {
			return nil, err
		}
		key, err := String(k)
		if err != nil {
			return nil, err
		}
		result[key] = str
	}
	return result, nil
}
