package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userSystemdEnable)
}

var userSystemdEnable = &Command{
	"user.systemd.enable",
	"Enable and start a systemd user service",
	[]string{
		"Username",
		"Name of the service to enable",
	},
	nil,
	"Enable syncthing\n  systemd.enable alice syncthing",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		name := args[1]
		if _, err := datastores.Users.Get(username); err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}

		var (
			outBuf bytes.Buffer
			outStr string
		)
		if err := helper.Exec(nil, &outBuf, "systemctl", "--machine="+username+"@.host", "--user", "is-enabled", name); err != nil {
			if !strings.Contains(outBuf.String(), "No such file or directory") {
				return nil, fmt.Sprintf("Could not check if %s is enabled: %s (%s)", name, err, outBuf.String()), nil, StatusError, "", ""
			}
		} else {
			outStr = strings.TrimSpace(outBuf.String())
		}
		if outStr == "enabled" {
			return nil, fmt.Sprintf("Service %s is already enabled", name), nil, StatusSuccess, "", ""
		}
		apply = func(out Output) bool {
			out.Runningf("Enabling and starting service %s", name)
			if err := helper.Exec(nil, nil, "systemctl", "--machine="+username+"@.host", "--user", "enable", "--now", name); err != nil {
				out.Errorf("Could not enable %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Enabled and started %s", name)
			return true
		}
		return nil, fmt.Sprintf("Need to enable service %s for user %s", name, username), apply, StatusInfo, "", ""
	},
	nil,
}
