package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(userVscodeExtensionInstall)
}

var userVscodeExtensionInstall = Command{
	"user.vscode.extension.install",
	"Install a VSCode extension",
	[]string{
		"Username",
		"Extension identifier",
	},
	nil,
	"VSCode with go\n  snap.install code classic\n  user.vscode.extension.install alice golang.go",
	func(args []string) (result []string, message string, apply Apply, status Status, before, after string) {
		username := args[0]
		id := args[1]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		ok, err := datastores.VSCodeExtensions.Installed(user, id)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", id), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) bool {
			out.Runningf("Installing %s for %s", id, username)
			if err := helper.ExecAs(user.User, nil, nil, "code", "--install-extension", id); err != nil {
				out.Errorf("Could not install %s for %s: %s", id, username, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s for %s", id, username)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s for %s", id, username), apply, StatusInfo, "", ""
	},
	datastores.VSCodeExtensions.Reset,
}
