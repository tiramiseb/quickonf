package modules

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("user-in-group", UserInGroup)
	Register("user-password", UserPassword)
}

// UserInGroup makes sure a user is in a group
func UserInGroup(in interface{}, out output.Output) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for user, group := range data {
		if Dryrun {
			out.Infof("Would add user %s to group %s", user, group)
			continue
		}
		if _, err := helper.ExecSudo(nil, "", "adduser", user, group); err != nil {
			return err
		}
		out.Successf("User %s is in group %s", user, group)
	}
	return nil
}

// UserPassword makes sure a user exists and has the given password
func UserPassword(in interface{}, out output.Output) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for user, password := range data {
		if _, err := helper.Exec(nil, "", "getent", "passwd", user); err != nil {
			if Dryrun {
				out.Infof("Would create user %s", user)
				continue
			}
			if _, err := helper.ExecSudo(nil, "", "useradd", "--create-home", "--shell", "/bin/bash", user); err != nil {
				return err
			}
			out.Successf("Created user %s", user)
		}
		if Dryrun {
			out.Infof("Would set password (not displayed) for user %s", user)
			continue
		}
		if _, err := helper.ExecSudo(nil, fmt.Sprintf("%s:%s", user, password), "chpasswd"); err != nil {
			return err
		}
		out.Successf("Password set for user %s", user)
	}
	return nil
}
