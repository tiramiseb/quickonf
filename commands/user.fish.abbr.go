package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userFishAbbr)
}

var userFishAbbr = &Command{
	"user.fish.abbr",
	"Create an abbreviation for the Fish shell",
	[]string{
		"Username",
		"Name of the abbreviation",
		"Command to bind to the abbreviation",
	},
	nil,
	"Git checkout abbreviation\n  user.fish.abbr alice gco \"git checkout\"",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		name := args[1]
		command := args[2]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		existing, err := datastores.Fish.Abbreviation(user, name)
		if err != nil {
			return nil, fmt.Sprintf("Could not check abbreviation %s: %s", name, err), nil, StatusError, "", ""
		}
		if existing == command {
			return nil, fmt.Sprintf("Abbreviation %s is already defined", name), nil, StatusSuccess, "", command
		}
		apply = func(out Output) (success bool) {
			if existing == "" {
				out.Runningf("Adding abbreviation %s", name)
			} else {
				out.Runningf("Modifying abbreviation %s", name)
			}
			if err := helper.ExecAs(user.User, nil, nil, "fish", "-c", fmt.Sprintf("abbr --add %s \"%s\"", name, command)); err != nil {
				out.Errorf("Could not define abbreviation %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Abbreviation %s defined", name)
			return true
		}
		var verb string
		if existing == "" {
			verb = "add"
		} else {
			verb = "modify"
		}
		return nil, fmt.Sprintf("Need to %s abbreviation %s", verb, name), apply, StatusInfo, existing, command
	},
	datastores.Fish.Reset,
}
