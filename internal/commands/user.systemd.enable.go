package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(userSystemdEnable)
}

var userSystemdEnable = Command{
	"user.systemd.enable",
	"Enable and start a systemd user service",
	[]string{
		"Username",
		"Name of the service to enable",
	},
	nil,
	"Enable syncthing\n  systemd.enable alice syncthing",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		name := args[1]
		if _, err := datastores.Users.Get(username); err != nil {
			return nil, err.Error(), nil, StatusError
		}

		var out bytes.Buffer
		if err := helper.Exec(nil, &out, "systemctl", "--machine="+username+"@.host", "--user", "is-enabled", name); err != nil {
			return nil, fmt.Sprintf("Could not check if %s is enabled: %s (%s)", name, err, out.String()), nil, StatusError
		}
		outS := strings.TrimSpace(out.String())
		if outS == "enabled" {
			return nil, fmt.Sprintf("Service %s is already enabled", name), nil, StatusSuccess
		} else if outS == "disabled" || strings.Contains(outS, "No such file or directory") {
			apply = &Apply{
				"user.systemd.enable",
				fmt.Sprintf("Will enable service %s for user %s", name, username),
				func(out Output) bool {
					out.Runningf("Enabling and starting service %s", name)
					if err := helper.Exec(nil, nil, "systemctl", "--machine="+username+"@.host", "--user", "enable", "--now", name); err != nil {
						out.Errorf("Could not enable %s: %s", name, helper.ExecErr(err))
						return false
					}
					out.Successf("Enabled and started %s", name)
					return true
				},
			}
			return nil, fmt.Sprintf("Need to enable service %s for user %s", name, username), apply, StatusInfo
		}
		return nil, outS, nil, StatusError
	},
	nil,
}
