package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("add-user-to-group", AddUserToGroup)
}

// AddUserToGroup adds a user to a group
func AddUserToGroup(in interface{}, out output.Output) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for user, group := range data {
		if Dryrun {
			out.Info("Would add user " + user + " to group " + group)
			continue
		}
		helper.ExecSudo("adduser", user, group)
		if err != nil {
			return err
		}
		out.Success("User " + user + " added to group " + group)
	}
	return nil
}
