package input

import (
	"errors"
)

func SliceString(input interface{}, store map[string]interface{}) ([]string, error) {
	data, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("Data is not a list")
	}
	result := make([]string, len(data))
	for i, v := range data {
		str, err := resolveString(v, store)
		if err != nil {
			return nil, err
		}
		result[i] = str
	}
	return result, nil
}
