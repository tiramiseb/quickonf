package commands

import (
	"fmt"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(dpkgInstall)
}

var dpkgInstall = Command{
	"dpkg.install",
	"Install a package using dpkg (does not check the file exists first)",
	[]string{"Absolute path to the package file"},
	nil,
	"Install that awesome package\n  dpkg.install <confdir>/my-awesome-stuff.deb",
	func(args []string) (result []string, msg string, apply Apply, status Status) {
		path := args[0]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError
		}

		apply = func(out Output) bool {
			out.Infof("Waiting for dpkg to be available to install %s", path)
			datastores.DpkgMutex.Lock()
			defer datastores.DpkgMutex.Unlock()
			out.Runningf("Installing %s", path)
			if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "dpkg", "--install", path); err != nil {
				out.Errorf("Could not install %s: %s", path, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s", path)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", path), apply, StatusInfo
	},
	datastores.DpkgPackages.Reset,
}
