package commands

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userFishVariable)
}

var userFishVariable = &Command{
	"user.fish.variable",
	"Set an universal variable for the Fish shell",
	[]string{
		"Username",
		"Name of the variable",
		"Content of the variable",
	},
	nil,
	"Tide time color\n  user.fish.variable alice tide_time_color 000000",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		name := args[1]
		content := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		existing, err := datastores.Fish.Variable(user, name)
		if err != nil {
			return nil, fmt.Sprintf("Could not check variable %s: %s", name, err), nil, StatusError, "", ""
		}
		if strings.Contains(content, " ") && !(content[0] == '\'' && content[len(content)-1] == '\'') {
			content = "'" + content + "'"
		}
		if existing == content {
			return nil, fmt.Sprintf("Variable %s is already defined", name), nil, StatusSuccess, "", content
		}
		apply = func(out Output) (success bool) {
			if existing == "" {
				out.Runningf("Adding variable %s", name)
			} else {
				out.Runningf("Modifying variable %s", name)
			}
			if err := helper.ExecAs(user.User, nil, nil, "fish", "-c", fmt.Sprintf("set --universal %s \"%s\"", name, content)); err != nil {
				out.Errorf("Could not set variable %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Variable %s defined", name)
			return true
		}
		var verb string
		if existing == "" {
			verb = "add"
		} else {
			verb = "modify"
		}
		return nil, fmt.Sprintf("Need to %s variable %s", verb, name), apply, StatusInfo, existing, content
	},
	datastores.Fish.Reset,
}
