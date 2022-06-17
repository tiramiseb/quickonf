package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(gitConfig)
}

var gitConfig = Command{
	"git.config",
	"Set a git configuration parameter in the system",
	[]string{
		`Name (in the form "x.y")`,
		"Value",
	},
	nil,
	"Git parameter\n  git.config branch.autosetuprebase always",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]
		value := args[1]
		current, err := datastores.GitConfig.Get(datastores.FakeUserForSystem, name)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value of %s: %s", name, err), nil, StatusError, "", ""
		}
		if value == current {
			return nil, fmt.Sprintf("%s is already set to %s", name, current), nil, StatusSuccess, value, value
		}

		apply = func(out Output) bool {
			out.Runningf("Setting %s to %s", name, value)
			if err := helper.Exec(nil, nil, "git", "config", "--system", "--add", name, value); err != nil {
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
