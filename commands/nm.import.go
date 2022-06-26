package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(nmImport)
}

var nmImport = &Command{
	"nm.import",
	"Import a VPN configuration in NetworkManager",
	[]string{
		"Type of the VPN",
		"Path of the configuration file to import",
	},
	nil,
	"Import my VPN configuration\n  nm.import openvpn /opt/openvpn.conf",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		connType := args[0]
		path := args[1]

		nameParts := strings.Split(filepath.Base(path), ".")
		var name string
		switch len(nameParts) {
		case 0:
			return nil, fmt.Sprintf("no file name in %s", path), nil, StatusError, "", ""
		case 1, 2:
			name = nameParts[0]
		default:
			name = strings.Join(nameParts[0:len(nameParts)-2], ".")
		}

		conn, ok, err := datastores.NetworkManagerConnections.Get(name)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			if conn.Type == "vpn" {
				return nil, fmt.Sprintf("%s is already configured", name), nil, StatusSuccess, conn.String(), ""
			}
			return nil, fmt.Sprintf("%s exists but is not of the right type", name), nil, StatusError, conn.String(), ""
		}

		apply = func(out Output) bool {
			out.Runningf("Importing %s", name)
			if err := helper.Exec(nil, nil, "nmcli", "connection", "import", "type", connType, "file", path); err != nil {
				out.Errorf("Could not import %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Imported %s", name)
			return true
		}
		return nil, fmt.Sprintf("Need to import %s configuration from %s", connType, path), apply, StatusInfo, "", ""
	},
	datastores.NetworkManagerConnections.Reset,
}
