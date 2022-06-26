package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(nmWifi)
}

var nmWifi = &Command{
	"nm.wifi",
	"Have knowledge of a wifi network",
	[]string{
		"SSID of the wifi network",
		"PSK of the wifi network (password)",
	},
	nil,
	"My own wifi network\n  nm.wifi mynetwork n0tSecureP4ssw0rd",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		ssid := args[0]
		psk := args[1]

		conn, ok, err := datastores.NetworkManagerConnections.Get(ssid)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			if conn.Type != "wifi" {
				return nil, fmt.Sprintf("%s exists but is not a wifi connection", ssid), nil, StatusError, conn.String(), ""
			}
			if conn.Mode != "infrastructure" {
				return nil, fmt.Sprintf("%s exists but is not an wifi client configuration", ssid), nil, StatusError, conn.String(), ""
			}
			if conn.PSK == psk {
				return nil, fmt.Sprintf("%s is already configured", ssid), nil, StatusSuccess, conn.String(), ""
			}
			msg = fmt.Sprintf("Need to change PSK for wifi network %s", ssid)
			before = conn.String()
			after = psk
			apply = func(out Output) bool {
				out.Infof("Changing PSK for wifi network %s", ssid)
				if err := helper.Exec(nil, nil, "nmcli", "connection", "modify", ssid, "802-11-wireless-security.psk", psk); err != nil {
					out.Errorf("Could not change PSK for wifi network %s: %s", ssid, helper.ExecErr(err))
					return false
				}
				out.Successf("Changed PSK for wifi network %s", ssid)
				return true
			}
		} else {
			msg = fmt.Sprintf("Need to store wifi network %s", ssid)
			apply = func(out Output) bool {
				out.Runningf("Storing wifi network %s", ssid)
				if err := helper.Exec(nil, nil, "nmcli", "connection", "add", "con-name", ssid, "type", "wifi", "ssid", ssid, "--", "802-11-wireless-security.key-mgmt", "wpa-psk", "802-11-wireless-security.psk", psk); err != nil {
					out.Errorf("Could not store wifi network %s: %s", ssid, helper.ExecErr(err))
					return false
				}
				out.Successf("Stored wifi network %s", ssid)
				return true
			}
		}
		return nil, msg, apply, StatusInfo, before, after
	},
	datastores.NetworkManagerConnections.Reset,
}
