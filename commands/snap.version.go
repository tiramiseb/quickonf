package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
)

func init() {
	register(snapVersion)
}

var snapVersion = Command{
	"snap.version",
	"Get the version of an installed snap package",
	[]string{
		"Name of the package",
	},
	[]string{
		"Version of the package",
	},
	"Example package version\n  version = snap.version obsidian",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]

		pkg, ok, err := datastores.Snap.Get(name)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if !ok {
			return []string{""}, fmt.Sprintf("Package %s is not installed", name), nil, StatusSuccess, "", ""
		}
		return []string{pkg.Version}, fmt.Sprintf("Package %s is installed in version %s", name, pkg.Version), nil, StatusSuccess, "", ""
	},
	nil,
}
