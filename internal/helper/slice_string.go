package helper

import (
	"errors"
)

// SliceString transforms a slice of interfaces to a slice of strings
func SliceString(input interface{}) ([]string, error) {
	data, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("Data is not a list")
	}
	result := make([]string, len(data))
	for i, v := range data {
		str, err := String(v)
		if err != nil {
			return nil, err
		}
		result[i] = str
	}
	return result, nil
}
