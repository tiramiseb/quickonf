package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
)

func init() {
	register(dpkgVersion)
}

var dpkgVersion = Command{
	"dpkg.version",
	"Get the version of an installed dpkg package",
	[]string{
		"Name of the package",
	},
	[]string{
		"Version of the package",
	},
	"Example package version\n  version = dpkg.version gnome-shell",
	func(args []string) (result []string, msg string, apply Apply, status Status) {
		name := args[0]

		pkg, ok, err := datastores.DpkgPackages.Get(name)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		if !ok {
			return []string{""}, fmt.Sprintf("Package %s is not installed", name), nil, StatusSuccess
		}
		return []string{pkg.Version}, fmt.Sprintf("Package %s is installed in version %s", name, pkg.Version), nil, StatusSuccess
	},
	nil,
}
