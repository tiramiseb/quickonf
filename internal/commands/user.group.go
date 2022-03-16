package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(userGroup)
}

var userGroup = Command{
	"user.group",
	"Add the user to the group (if the group does not exist, it is created)",
	[]string{
		"Username",
		"Groupname",
	},
	nil,
	"Allow alice to dial with a modem\n  user.group alice dialout",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		groupname := args[1]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if user.Group.Name == groupname {
			return nil, fmt.Sprintf("%s is already a member of %s", username, groupname), nil, StatusSuccess, "", ""
		}
		for _, g := range user.Groups {
			if g.Name == groupname {
				return nil, fmt.Sprintf("%s is already a member of %s", username, groupname), nil, StatusSuccess, "", ""
			}
		}

		_, exists, err := datastores.Users.GetGroup(groupname)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}

		var needMessage string
		if exists {
			needMessage = fmt.Sprintf("Need to create group %s and make %s a member of it", groupname, username)
		} else {
			needMessage = fmt.Sprintf("Need to make %s a member of group %s", username, groupname)
		}

		apply = func(out Output) bool {
			if !exists {
				out.Runningf("Creating group %s", groupname)
				if err := helper.Exec(nil, nil, "groupadd", groupname); err != nil {
					out.Errorf("Could not create group %s: %s", groupname, err)
					return false
				}
				out.Infof("Created group %s", groupname)
			}
			out.Runningf("Adding %s to group %s", username, groupname)
			if err := helper.Exec(nil, nil, "usermod", "--append", "--groups", groupname); err != nil {
				out.Errorf("Could not add %s to group %s: %s", username, groupname, err)
				return false
			}
			out.Successf("User %s added to group %s", username, groupname)
			return true
		}
		return nil, needMessage, apply, StatusInfo, "", ""
	},
	datastores.Users.Reset,
}
