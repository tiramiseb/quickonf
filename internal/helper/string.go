package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var filepathRe = regexp.MustCompile("<file:([^>]*)>")

// String returns a string from an interface...
// ... replacing occurrences of "<file:xxx>" to the content of the given file if it exists (path relative to the path of the configuration file).
// ... replacing occurrences of "<xxx>" with the value of xxx in the store.
func String(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	var str string
	switch val := v.(type) {
	case int:
		str = strconv.Itoa(val)
	case string:
		str = val
	default:
		return "", fmt.Errorf(`value "%v" is not a string`, v)
	}
	str = replaceFile(str)
	str = replaceStore(str)
	return str, nil
}

// replaceFile replaces values from files contents
func replaceFile(str string) string {
	return filepathRe.ReplaceAllStringFunc(str, func(src string) string {
		subpath := filepathRe.FindStringSubmatch(src)[1]
		var fpath string
		if subpath[0] == '/' {
			fpath = subpath
		} else {
			fpath = filepath.Join(Basepath, subpath)
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			return str
		}
		return string(data)
	})
}

// replaceStore replaces values from the store
func replaceStore(str string) string {
	for key, val := range store {
		str = strings.ReplaceAll(str, "<"+key+">", val)
	}
	return str
}
