package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userGitConfig)
}

var userGitConfig = &Command{
	"user.git.config",
	"Set a git configuration parameter for a user",
	[]string{
		"Username",
		`Name (in the form "x.y")`,
		"Value",
	},
	nil,
	"Git parameter\n  user.git.config alice branch.autosetuprebase always",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		name := args[1]
		value := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		current, err := datastores.GitConfig.Get(user, name)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", name, err), nil, StatusError, "", ""
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", name, current), nil, StatusSuccess, current, value
		}

		apply = func(out Output) bool {
			out.Runningf("Setting %s to %s", name, value)
			if err := helper.ExecAs(user.User, nil, nil, "git", "config", "--global", "--add", name, value); err != nil {
				out.Errorf("Could not set %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Set %s to %s", name, value)
			return true
		}
		return nil, fmt.Sprintf("Need to set %s to %s", name, value), apply, StatusInfo, current, value
	},
	func() {
		datastores.GitConfig.Reset()
		datastores.Users.Reset()
	},
}
