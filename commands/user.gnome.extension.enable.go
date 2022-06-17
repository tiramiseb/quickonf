package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(gnomeExtensionEnable)
}

var gnomeExtensionEnable = Command{
	"user.gnome.extension.enable",
	"Enable a GNOME Shell extension",
	[]string{
		"Username",
		"UUID of the extension",
	},
	nil,
	"Dash to dock\n  gnome.extension.install dash-to-dock@micxgx.gmail.com\n  user.gnome.extension.enable alice dash-to-dock@micxgx.gmail.com",
	func(args []string) (result []string, message string, apply Apply, status Status, before, after string) {
		username := args[0]
		uuid := args[1]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		ok, err := datastores.GnomeExtensions.Enabled(user, uuid)
		if err != nil {
			return nil, fmt.Sprintf("Could not check if %s is enabled: %s", uuid, err), nil, StatusError, "", ""
		}
		if ok {
			return nil, fmt.Sprintf("%s is already enabled", uuid), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) (success bool) {
			out.Runningf("Enabling %s", uuid)
			if err := helper.Exec(nil, nil, "gnome-shell-extension-tool", "--enable-extension", uuid); err != nil {
				out.Errorf("Could not enable %s: %s", uuid, err)
				return false
			}
			out.Successf("Enabled %s", uuid)
			return true
		}

		return nil, fmt.Sprintf("Need to enable %s", uuid), nil, StatusInfo, "", ""
	},
	datastores.GnomeExtensions.Reset,
}
