package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("nm-import-openvpn", NetworkManagerImport)
	Register("nm-wifi", NetworkManagerWifi)
}

// NetworkManagerWifi sets the known wifi networks with WPA-PSK security
func NetworkManagerWifi(in interface{}, out output.Output) error {
	out.InstructionTitle("Network manager wifi")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for net, ssid := range data {
		ok, err := nmConnectionExists(net)
		if err != nil {
			return err
		}
		if ok {
			out.Infof("Connection %s already exists", net)
			continue
		}
		if Dryrun {
			out.Infof("Would add wifi connection %s", net)
			continue
		}
		if _, err := helper.Exec(nil, "nmcli", "connection", "add", "con-name", net, "type", "wifi", "ssid", net, "--", "802-11-wireless-security.key-mgmt", "wpa-psk", "802-11-wireless-security.psk", ssid); err != nil {
			return err
		}
		out.Successf("Added wifi connection %s", net)
	}
	return nil
}

// NetworkManagerImport imports a configuration into network manager
func NetworkManagerImport(in interface{}, out output.Output) error {
	out.InstructionTitle("Importing configuration into network manager")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				out.Alertf("File %s does not exist", path)
				continue
			}
			return err
		}
		nameParts := strings.Split(filepath.Base(path), ".")
		var name string
		switch len(nameParts) {
		case 0:
			return fmt.Errorf("No file name in %s", path)
		case 1, 2:
			name = nameParts[0]
		default:
			name = strings.Join(nameParts[0:len(nameParts)-2], ".")
		}
		path = helper.Path(path)

		ok, err := nmConnectionExists(name)
		if err != nil {
			return err
		}
		if ok {
			out.Infof("%s is already configured", name)
			continue
		}

		if Dryrun {
			out.Info("Would import " + path)
			continue
		}

		if _, err := helper.Exec(nil, "nmcli", "connection", "import", "type", "openvpn", "file", path); err != nil {
			return err
		}
	}
	return nil
}

func nmConnectionExists(name string) (bool, error) {
	shownB, err := helper.Exec(nil, "nmcli", "connection", "show")
	if err != nil {
		return false, err
	}
	shown := string(shownB)
	for _, line := range strings.Split(shown, "\n") {
		if strings.HasPrefix(line, name) {
			return true, nil
		}
	}
	return false, nil
}
