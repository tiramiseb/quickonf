package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(gitUserConfig)
}

var gitUserConfig = Command{
	"git.user.config",
	"Set a git configuration parameter for a user",
	[]string{
		"Username",
		`Name (in the form "x.y")`,
		"Value",
	},
	nil,
	"Git parameter\n  git.config branch.autosetuprebase always", func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		name := args[1]
		value := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError
		}
		current, err := datastores.GitConfig.Get(user, name)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", name, err), nil, StatusError
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", name, current), nil, StatusSuccess
		}

		apply = &Apply{
			"git.user.config",
			fmt.Sprintf("Will set %s to %s", name, value),
			func(out Output) bool {
				out.Infof("Setting %s to %s", name, value)
				if err := helper.ExecAs(user.User, nil, nil, "git", "config", "--global", "--add", name, value); err != nil {
					out.Errorf("Could not set %s: %s", name, helper.ExecErr(err))
					return false
				}
				out.Successf("Set %s to %s", name, value)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to set %s to %s", name, value), apply, StatusInfo
	},
	func() {
		datastores.Dconf.Reset()
		datastores.Users.Reset()
	},
}
