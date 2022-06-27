package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(dpkgReconfigure)
}

var dpkgReconfigure = &Command{
	"dpkg.reconfigure",
	"Reconfigure an installed dpkg package",
	[]string{
		"Name of the package",
	},
	nil,
	"libdvd-pkg\n  debconf.set libdvd-pkg libdvd-pkg/first-install .\n  debconf.set libdvd-pkg libdvd-pkg/post-invoke_hook-install true\n  apt.install libdvd-pkg\n  dpkg.reconfigure libdvd-pkg",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]
		apply = func(out Output) bool {
			out.Infof("Waiting for dpkg to be available to reconfigure %s", name)
			datastores.DpkgMutex.Lock()
			defer datastores.DpkgMutex.Unlock()
			out.Runningf("Reconfiguring %s", name)
			// _, err := helper.ExecSudo(nil, "", "dpkg-reconfigure", "--frontend", "noninteractive", name)

			if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "dpkg-reconfigure", name); err != nil {
				out.Errorf("Could not reconfigure %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Reconfigured %s", name)
			return true

		}
		return nil, fmt.Sprintf("Need to reconfigure %s", name), apply, StatusInfo, "", ""
	},
	nil,
}
