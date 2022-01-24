package commands

import (
	"fmt"
	"regexp"
)

func init() {
	register(regexpSubstring)
}

var regexpSubstring = Command{
	"regexp.submatch",
	"Find submatches using a regexp",
	[]string{
		"Regexp",
		"Source string",
	},
	[]string{
		"First submatch",
		"Second submatch",
		"...",
	},
	"Find src\n  webpage = http.get.var http://www.example.com\n  src = regexp.substring \"<script .*src=(.*)>\" <webpage>\n  ...",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		reg := args[0]
		source := args[1]

		re, err := regexp.Compile(reg)
		if err != nil {
			return nil, fmt.Sprintf("%s is not a valid regexp: %s", reg, err), nil, StatusError
		}

		results := re.FindStringSubmatch(source)
		return results[1:], fmt.Sprintf("Matched regexp %v", results[1:]), nil, StatusSuccess
	},
	nil,
}
