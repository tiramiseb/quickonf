package commands

import (
	"fmt"
	"regexp"
	"strings"
)

func init() {
	register(regexpSubmatch)
}

var regexpSubmatch = &Command{
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
	"Find src\n  webpage = http.get.var http://www.example.com\n  src = regexp.submatch \"<script .*src=(.*)>\" <webpage>\n  ...",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		reg := args[0]
		source := args[1]

		re, err := regexp.Compile(reg)
		if err != nil {
			return nil, fmt.Sprintf("%s is not a valid regexp: %s", reg, err), nil, StatusError, "", ""
		}

		results := re.FindStringSubmatch(source)
		if results == nil {
			return nil, fmt.Sprintf("No match for regexp %s", reg), nil, StatusError, "", ""
		}
		return results[1:], "Matched regexp", nil, StatusSuccess, "", `"` + strings.Join(results[1:], `", "`) + `"`
	},
	nil,
}
