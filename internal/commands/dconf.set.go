package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(dconfSet)
}

var dconfSet = Command{
	"dconf.set",
	"Set a dconf parameter",
	[]string{
		"Username",
		`Key (in the form "/x/y/z")`,
		"Value",
	},
	nil,
	"Show date in GNOME\n  dconf.set /org/gnome/desktop/interface/clock-show-date true", func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		key := args[1]
		value := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError
		}
		current, err := datastores.Dconf.Get(user, key)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", key, err), nil, StatusError
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", key, current), nil, StatusSuccess
		}

		// Convert value to something understandable by dconf
		_, errBool := strconv.ParseBool(value)
		_, errFloat := strconv.ParseFloat(value, 64)
		_, errInt := strconv.ParseInt(value, 10, 64)
		if errBool != nil && errFloat != nil && errInt != nil {
			// Not parseable as bool, float or int...
			if value == "[]" {
				value = "@as []"
			} else {
				value = "\"" + strings.ReplaceAll(value, "\"", "\\\"") + "\""
			}
		}

		apply = &Apply{
			"dconf.set",
			fmt.Sprintf("Will set %s to %s", key, value),
			func(out Output) bool {
				out.Infof("Setting %s to %s", key, value)
				wait, err := helper.ExecAs(user, nil, nil, "dconf", "write", key, value)
				if err != nil {
					out.Errorf("Could not set %s: %s", key, err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not set %s: %s", key, err)
					return false
				}
				out.Successf("Set %s to %s", key, value)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to set %s to %s", key, value), apply, StatusInfo
	},
	nil,
}
