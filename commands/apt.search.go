package commands

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(aptSearch)
}

var aptSearch = &Command{
	"apt.search",
	"Search for a package according to a regexp",
	[]string{"Regexp to match the package name"},
	[]string{"Name of the found package (last one in aphabetical order, if there are multiple matching packages)"},
	"Install the latest NVidia driver\n  pkg = apt.search ^nvidia-driver-.*-open$\n  apt.install <pkg>",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		pattern := args[0]

		var out strings.Builder
		if err := helper.Exec(nil, &out, "apt-cache", "--names-only", "search", pattern); err != nil {
			return nil, fmt.Sprintf("Could not search for package matching %q: %s", pattern, err), nil, StatusError, "", ""
		}
		lines := strings.Split(out.String(), "\n")
		if len(lines) == 0 {
			return nil, fmt.Sprintf("No package matches %q", pattern), nil, StatusError, "", ""
		}
		fields := strings.Fields(lines[len(lines)-1])
		if len(fields) == 0 || fields[0] == "" {
			return nil, fmt.Sprintf("No package matches %q", pattern), nil, StatusError, "", ""
		}
		pkgName := fields[0]

		result = []string{pkgName}
		return result, fmt.Sprintf("Package %q matches %q", pkgName, pattern), nil, StatusSuccess, "", pkgName
	},
	nil,
}
