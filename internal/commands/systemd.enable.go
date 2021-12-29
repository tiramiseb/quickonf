package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(systemdEnable)
}

var systemdEnable = Command{
	"systemd.enable",
	"Enable and start a systemd service",
	[]string{
		"Name of the service to enable",
	},
	nil,
	"Enable nextdns\n  systemd.enable nextdns",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		name := args[0]

		var out bytes.Buffer
		wait, err := helper.Exec(nil, &out, "systemctl", "is-enabled", name)
		if err != nil {
			return nil, fmt.Sprintf("Could not check if %s is enabled: %s", name, err), nil, StatusError
		}
		if err := wait(); err != nil && err.Error() != "exit status 1" {
			return nil, fmt.Sprintf("Could not check if %s is enabled: %s", name, err), nil, StatusError
		}
		outS := strings.TrimSpace(out.String())
		if outS == "enabled" {
			return nil, fmt.Sprintf("Service %s is already enabled", name), nil, StatusSuccess
		} else if outS == "disabled" || strings.Contains(outS, "No such file or directory") {
			apply = &Apply{
				"systemd.enable",
				fmt.Sprintf("Will enable service %s", name),
				func(out Output) bool {
					out.Infof("Enabling and starting service %s", name)
					wait, err := helper.Exec(nil, nil, "systemctl", "enable", "--now", name)
					if err != nil {
						out.Errorf("Could not enable %s: %s", name, err)
						return false
					}
					if err := wait(); err != nil {
						out.Errorf("Could not enable %s: %s", name, err)
						return false
					}
					out.Successf("Enabled and started %s", name)
					return true
				},
			}
			return nil, fmt.Sprintf("Need to enable service %s", name), apply, StatusInfo
		}
		return nil, outS, nil, StatusError
	},
	nil,
}
