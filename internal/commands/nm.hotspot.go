package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(nmHotspot)
}

// TODO Change autoconnect (needs changes in nm datastore)
var nmHotspot = Command{
	"nm.hotspot",
	"Setup a wifi access point (hotspot) - only one can be configured",
	[]string{
		"Network Interface name",
		"SSID of the hotspot",
		"PSK of the hotspot (password)",
		"Autoconnect (one of: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)",
	},
	nil,
	"Share connection\n  nm.hotspot wlp4s0 mynetwork n0tSecureP4ssw0rd true",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		iface := args[0]
		ssid := args[1]
		psk := args[2]
		autoconnect, err := strconv.ParseBool(args[3])
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		autoconnectYesno := "no"
		if autoconnect {
			autoconnectYesno = "yes"
		}

		if ok, err := helper.NetworkInterfaceExists(iface); err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		} else if !ok {
			return nil, fmt.Sprintf("interface %s does not exist", iface), nil, StatusError, "", ""
		}

		conn, ok, err := datastores.NetworkManagerConnections.Get("Hotspot")
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			if conn.Type != "wifi" || conn.Mode != "ap" {
				return nil, "Connection with name Hotspot exists but is not a hotspot", nil, StatusError, conn.String(), ""
			}
			var fieldsToChange []string
			if conn.InterfaceName == iface && conn.SSID == ssid && conn.PSK == psk {
				return nil, fmt.Sprintf("%s is already configured", ssid), nil, StatusSuccess, conn.String(), ""
			}
			if conn.InterfaceName != iface {
				fieldsToChange = append(fieldsToChange, "interface")
			}
			if conn.SSID != ssid {
				fieldsToChange = append(fieldsToChange, "SSID")
			}
			if conn.PSK != psk {
				fieldsToChange = append(fieldsToChange, "PSK")
			}
			msg = fmt.Sprintf("Need to change %s for hotspot", strings.Join(fieldsToChange, ", "))
			before = conn.String()
			after = "Interface: " + iface + "\nAutoconnect: " + autoconnectYesno + "\nSSID: " + ssid + "\nPSK: " + psk
			apply = func(out Output) bool {
				if conn.InterfaceName != iface {
					out.Info("Changing interface for hotspot")
					if err := helper.Exec(nil, nil, "nmcli", "connection", "modify", "Hotspot", "connection.interface-name", iface); err != nil {
						out.Errorf("Could not change SSID for hotspot: %s", helper.ExecErr(err))
						return false
					}
				}
				if conn.SSID != ssid {
					out.Info("Changing SSID for hotspot")
					if err := helper.Exec(nil, nil, "nmcli", "connection", "modify", "Hotspot", "802-11-wireless.ssid", ssid); err != nil {
						out.Errorf("Could not change SSID for hotspot: %s", helper.ExecErr(err))
						return false
					}
				}
				if conn.PSK != psk {
					out.Info("Changing PSK for hotspot")
					if err := helper.Exec(nil, nil, "nmcli", "connection", "modify", "Hotspot", "802-11-wireless-security.psk", psk); err != nil {
						out.Errorf("Could not change PSK for hotspot: %s", helper.ExecErr(err))
						return false
					}
				}
				out.Successf("Changed %s for hotspot", strings.Join(fieldsToChange, ", "))
				return true
			}
		} else {
			msg = "Need to store hotspot"
			apply = func(out Output) bool {
				out.Running("Storing hotspot")
				if err := helper.Exec(
					nil, nil, "nmcli", "connection", "add",
					"con-name", "Hotspot",
					"type", "wifi",
					"ifname", iface,
					"ssid", ssid,
					"mode", "ap",
					"autoconnect", autoconnectYesno,
					"--",
					"connection.autoconnect", autoconnectYesno,
					"802-11-wireless-security.key-mgmt", "wpa-psk",
					"802-11-wireless-security.psk", psk,
					"ipv4.method", "shared",
				); err != nil {
					out.Errorf("Could not store hotspot: %s", helper.ExecErr(err))
					return false
				}
				out.Success("Stored hotspot")
				return true
			}
		}
		return nil, msg, apply, StatusInfo, before, after
	},
	datastores.NetworkManagerConnections.Reset,
}
