package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
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
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		pkg := args[0]
		_, ok, err := datastores.DpkgPackages.Get(pkg)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", pkg), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) bool {
			out.Infof("Waiting for dpkg to be available to install %s", pkg)
			datastores.DpkgMutex.Lock()
			defer datastores.DpkgMutex.Unlock()
			out.Runningf("Installing %s", pkg)
			if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "install", pkg); err != nil {
				out.Errorf("Could not install %s: %s", pkg, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s", pkg)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", pkg), apply, StatusInfo, "", ""
	},
	datastores.DpkgPackages.Reset,
}
