package commands

import (
	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(aptUpgrade)
}

var aptUpgrade = Command{
	"apt.upgrade",
	"Upgrade packages from apt repositories",
	nil,
	nil,
	"Upgrade packages\n  apt.upgrade",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		apply = &Apply{
			"apt.upgrade",
			"Will upgrade apt packages",
			func(out Output) bool {
				out.Info("Waiting for dpkg to be available to upgrade packages")
				datastores.DpkgMutex.Lock()
				defer datastores.DpkgMutex.Unlock()
				out.Infof("Updating packages list")
				if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "update"); err != nil {
					out.Errorf("Could not update packages list: %s", helper.ExecErr(err))
					return false
				}
				out.Infof("Upgrading packages")
				if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "upgrade"); err != nil {
					out.Errorf("Could not upgrade packages: %s", helper.ExecErr(err))
					return false
				}
				out.Successf("Upgraded all packages")
				return true
			},
		}
		return nil, "Need to upgrade APT packages", apply, StatusInfo
	},
	nil,
}
