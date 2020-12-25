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
			out.Infof("Would add user %s to group %s", user, group)
			continue
		}
		if _, err := helper.ExecSudo(nil, "adduser", user, group); err != nil {
			return err
		}
		out.Successf("User %s added to group %s", user, group)
	}
	return nil
}
