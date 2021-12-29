package commands

import (
	"fmt"
	"os"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(debconfSet)
}

var debconfSet = Command{
	"debconf.set",
	"Set a debconf parameter",
	[]string{
		"Package where the parameter belong",
		"Name of the parameter",
		"Value to apply to the parameter",
	},
	nil,
	"Install MS fonts\n  debconf.set ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula true\n  apt.install ttf-mscorefonts-installer",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		pkg := args[0]
		name := args[1]
		value := args[2]
		param, ok, err := datastores.Debconf.Get(pkg, name)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		var verb string
		if ok {
			if param.Value == value {
				return nil, fmt.Sprintf("%s already has value %s", name, value), nil, StatusSuccess
			}
			verb = "change"
		} else {
			verb = "set"
		}
		apply = &Apply{
			"debconf.set",
			fmt.Sprintf("Will %s %s to %s", verb, name, value),
			func(out Output) bool {
				out.Infof("Setting %s to %s", name, value)
				tmpfile, err := os.CreateTemp("", "quickonf-debconf")
				if err != nil {
					out.Errorf("Could not create temporary file: %s", err)
					return false
				}
				defer os.Remove(tmpfile.Name())
				if _, err := tmpfile.WriteString(fmt.Sprintf("%s %s select %s", pkg, name, value)); err != nil {
					tmpfile.Close()
					out.Errorf("Could not write to temporary file: %s", err)
					return false
				}
				if err := tmpfile.Close(); err != nil {
					out.Errorf("Could not close temporary file: %s", err)
					return false
				}
				out.Infof("Waiting for dpkg to be available to set debconf value %s", name)
				datastores.DpkgMutex.Lock()
				defer datastores.DpkgMutex.Unlock()
				wait, err := helper.Exec(nil, nil, "debconf-set-selections", tmpfile.Name())
				if err != nil {
					out.Errorf("Could not execute debconf-set-selections: %s", err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not execute debconf-set-selections: %s", err)
					return false
				}
				out.Successf("%s set to %s", name, value)
				return true

			},
		}
		return nil, fmt.Sprintf("Need to %s %s to %s", verb, name, value), apply, StatusInfo
	},
	datastores.Debconf.Reset,
}
