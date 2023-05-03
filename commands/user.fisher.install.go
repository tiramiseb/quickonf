package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

const fisherFunctionURI = "https://raw.githubusercontent.com/jorgebucaran/fisher/main/functions/fisher.fish"

func init() {
	register(userFisherInstall)
}

var userFisherInstall = &Command{
	"user.fisher.install",
	"Install a fish plugin using Fisher",
	[]string{
		"Username",
		"Name/URI of the plugin",
	},
	nil,
	"Fish abbreviation tips\n  user.fisher.install alice gazorby/fish-abbreviation-tips",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		plugin := args[1]
		user, err := datastores.Users.Get(username)
		if err != nil {
			return nil, fmt.Sprintf("Could not get user %s: %s", username, err), nil, StatusError, "", ""
		}
		installed, err := datastores.Fish.IsPluginInstalled(user, plugin)
		if err != nil {
			return nil, fmt.Sprintf("Could not check plugin %s: %s", plugin, err), nil, StatusError, "", ""
		}
		if installed {
			return nil, fmt.Sprintf("Plugin %s is already installed", plugin), nil, StatusSuccess, "", ""
		}

		apply = func(out Output) (success bool) {
			out.Running("Checking fisher")
			var mustInstallFisher bool
			if err := helper.ExecAs(user.User, nil, nil, "fish", "-c", "fisher -v"); err != nil {
				if helper.ExecErrCode(err) == 127 {
					mustInstallFisher = true
				} else {
					out.Errorf("Could not check fisher: %s", helper.ExecErr(err))
					return false
				}
			}
			if mustInstallFisher {
				out.Running("Installing fisher")
				resp, err := http.Get(fisherFunctionURI)
				if err != nil {
					out.Errorf("Could not download fisher installer: %s", err)
					return false
				}
				f, err := os.CreateTemp("", "quickonf-fisher-*")
				if err != nil {
					resp.Body.Close()
					out.Errorf("Could not create temporary file for fisher installer: %s", err)
					return false
				}
				_, err = io.Copy(f, resp.Body)
				resp.Body.Close()
				if err != nil {
					out.Errorf("Could not write fisher installer to temporary file: %s", err)
					return false
				}

				if err := helper.ExecAs(user.User, nil, nil, "fish", "-c", fmt.Sprintf("source %s && fisher install jorgebucaran/fisher", f.Name())); err != nil {
					out.Errorf("Could not install fisher: %s", helper.ExecErr(err))
					return false
				}
			}
			if err := helper.ExecAs(user.User, nil, nil, "fish", "-c", fmt.Sprintf("fisher install %s", plugin)); err != nil {
				out.Errorf("Could not install %s: %s", plugin, helper.ExecErr(err))
				return false
			}
			out.Successf("Plugin %s installed", plugin)
			return true
		}
		return nil, fmt.Sprintf("Need to install fish plugin %s", plugin), apply, StatusInfo, "", ""
	},
	datastores.Fish.Reset,
}
