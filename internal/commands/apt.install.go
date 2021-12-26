package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
	"github.com/tiramiseb/quickonf/internal/commands/shared"
)

func init() {
	register(apt)
}

var apt = Command{
	"apt.install",
	"Install a package using apt",
	[]string{"Name of the package to install"},
	nil,
	"Install the \"ipcalc\" tool\n  apt.install ipcalc",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		pkg := args[0]
		ok, err := shared.DpkgPackages.Installed(pkg)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", pkg), nil, StatusSuccess
		}

		apply = &Apply{
			"apt.install",
			fmt.Sprintf("Will install %s", pkg),
			func(out Output) bool {
				out.Infof("Waiting for apt to be available to install %s", pkg)
				shared.DpkgMutex.Lock()
				defer shared.DpkgMutex.Unlock()
				wait, err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "install", pkg)
				if err != nil {
					out.Errorf("Could not install %s: %s", pkg, err)
					return false
				}
				out.Infof("Installing %s", pkg)
				if err := wait(); err != nil {
					out.Errorf("Could not install %s: %s", pkg, err)
					return false
				}
				out.Successf("Installed %s", pkg)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to install %s", pkg), apply, StatusInfo
	},
	shared.DpkgPackages.Reset,
}
