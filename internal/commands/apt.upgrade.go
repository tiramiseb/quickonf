package commands

import "github.com/tiramiseb/quickonf/internal/commands/helper"

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
				out.Infof("Updating packages list")
				wait, err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "update")
				if err != nil {
					out.Errorf("Could not update packages list: %s", err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not update packages list: %s", err)
					return false
				}
				out.Infof("Upgrading packages")
				wait, err = helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "upgrade")
				if err != nil {
					out.Errorf("Could not upgrade packages: %s", err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not upgrade packages: %s", err)
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
