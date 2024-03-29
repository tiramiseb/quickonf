package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(flatpakInstall)
}

var flatpakInstall = &Command{
	"flatpak.install",
	"Install a package using flatpak",
	[]string{"Application ID of the package to install"},
	nil,
	"Install \"anydesk\"\n  flatpak.install com.anydesk.Anydesk",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		appID := args[0]
		_, ok, err := datastores.Flatpak.Get(appID)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", appID), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) bool {
			out.Runningf("Installing %s", appID)
			if err := helper.Exec(nil, nil, "flatpak", "install", "--noninteractive", "--assumeyes", appID); err != nil {
				out.Errorf("Could not install %s: %s", appID, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s", appID)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", appID), apply, StatusInfo, "", ""
	},
	datastores.Flatpak.Reset,
}
