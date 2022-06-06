package commands

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(userGoEnv)
}

var userGoEnv = Command{
	"user.go.env",
	"Set a Go environment variable",
	[]string{
		"Username",
		"Variable name (automatically changed to uppercase)",
		"Value",
	},
	nil,
	"Go path\n  user.go.env alice gopath /home/alice/GO",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		variable := strings.ToUpper(args[1])
		value := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		current, err := datastores.GoEnv.Get(user, variable)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", variable, err), nil, StatusError, "", ""
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", variable, current), nil, StatusSuccess, current, value
		}

		apply = func(out Output) bool {
			out.Runningf("Setting %s to %s", variable, value)
			if err := helper.ExecAs(user.User, nil, nil, "go", "env", "-w", fmt.Sprintf(`%s="%s"`, variable, value)); err != nil {
				out.Errorf("Could not set %s: %s", variable, helper.ExecErr(err))
				return false
			}
			out.Successf("Set %s to %s", variable, value)
			return true
		}
		return nil, fmt.Sprintf("Need to set %s", variable), apply, StatusInfo, current, value
	},
	func() {
		datastores.GoEnv.Reset()
		datastores.Users.Reset()
	},
}
