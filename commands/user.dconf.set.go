package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userDconfSet)
}

var userDconfSet = &Command{
	"user.dconf.set",
	"Set a dconf parameter",
	[]string{
		"Username",
		`Key (in the form "/x/y/z")`,
		"Value",
	},
	nil,
	"Show date in GNOME\n  user.dconf.set /org/gnome/desktop/interface/clock-show-date true",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		key := args[1]
		value := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		current, err := datastores.Dconf.Get(user, key)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", key, err), nil, StatusError, "", ""
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", key, current), nil, StatusSuccess, current, value
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

		apply = func(out Output) bool {
			out.Runningf("Setting %s to %s", key, value)
			env := []string{fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/%d/bus", user.Uid)}
			if err := helper.ExecAs(user.User, env, nil, "dconf", "write", key, value); err != nil {
				out.Errorf("Could not set %s: %s", key, helper.ExecErr(err))
				return false
			}
			out.Successf("Set %s to %s", key, value)
			return true
		}
		return nil, fmt.Sprintf("Need to set %s to %s", key, value), apply, StatusInfo, current, value
	},
	func() {
		datastores.Dconf.Reset()
		datastores.Users.Reset()
	},
}
