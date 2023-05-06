package commands

import (
	"github.com/tidwall/gjson"
)

func init() {
	register(jsonGet)
}

var jsonGet = &Command{
	"json.get",
	"Get a JSON value from a variable",
	[]string{
		"Variable containing the JSON object",
		"Path of the value to get",
	},
	[]string{
		"Value",
	},
	"Get a value\n  val = json.get \"{\\\"download\\\":{\\\"url\\\":\\\"http://www.example.com/\\\"}}\" download.url",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		data := args[0]
		path := args[1]
		if ok := gjson.Valid(data); !ok {
			return []string{""}, "First argument is not a JSON object", nil, StatusError, "", ""
		}
		value := gjson.Get(data, path)
		str := value.String()
		return []string{str}, "Read value from JSON", nil, StatusSuccess, "", str
	},
	nil,
}
