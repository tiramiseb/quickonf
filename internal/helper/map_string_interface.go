package helper

import (
	"errors"
)

// MapStringInterface transforms an interface to a map from strings to interfaces
func MapStringInterface(input interface{}) (map[string]interface{}, error) {
	data, ok := input.(map[string]interface{})
	if !ok {
		return nil, errors.New("module configuration must be a map")
	}
	result := make(map[string]interface{}, len(data))
	for k, v := range data {
		key, err := String(k)
		if err != nil {
			return nil, err
		}
		result[key] = v
	}
	return result, nil
}
