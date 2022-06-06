package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(pipInstall)
}

var pipInstall = Command{
	"pip.install",
	"Install a Python package using pip3",
	[]string{"Name of the package to install"},
	nil,
	"Pip Passlib\n  pip.install passlib",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		pkg := args[0]
		pkgObj, ok, err := datastores.PipPackages.Get(pkg)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", pkg), nil, StatusSuccess, "Version " + pkgObj.Version, ""
		}

		apply = func(out Output) bool {
			out.Runningf("Installing %s", pkg)
			if err := helper.Exec(nil, nil, "pip3", "install", pkg); err != nil {
				out.Errorf("Could not install %s: %s", pkg, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s", pkg)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", pkg), apply, StatusInfo, "", ""
	},
	datastores.PipPackages.Reset,
}
