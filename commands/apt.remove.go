package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(aptRemove)
}

var aptRemove = Command{
	"apt.remove",
	"Remove a package using apt",
	[]string{"Name of the package to remove"},
	nil,
	"Remove the \"ipcalc\" tool\n  apt.remove ipcalc",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		pkg := args[0]
		_, ok, err := datastores.DpkgPackages.Get(pkg)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if !ok {
			return nil, fmt.Sprintf("%s is already not installed", pkg), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) bool {
			out.Infof("Waiting for dpkg to be available to remove %s", pkg)
			datastores.DpkgMutex.Lock()
			defer datastores.DpkgMutex.Unlock()
			out.Runningf("Removing %s", pkg)
			if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "remove", pkg); err != nil {
				out.Errorf("Could not remove %s: %s", pkg, helper.ExecErr(err))
				return false
			}
			out.Successf("Removed %s", pkg)
			return true
		}
		return nil, fmt.Sprintf("Need to remove %s", pkg), apply, StatusInfo, "", ""
	},
	datastores.DpkgPackages.Reset,
}
