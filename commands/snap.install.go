package commands

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/tiramiseb/quickonf/commands/datastores"
	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(snapInstall)
}

var snapInstall = Command{
	"snap.install",
	"Install a package using snap",
	[]string{
		"Name of the package to install",
		"Options, comma-separated, at least one of: stable, candidate, edge, beta, classic, dangerous",
	},
	nil,
	"Install node\n  snap.install node stable",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]
		options := strings.FieldsFunc(args[1], func(c rune) bool {
			return c == ',' || unicode.IsSpace(c)
		})
		var (
			channel   = "stable"
			classic   bool
			dangerous bool
			devmode   bool
		)
		for _, o := range options {
			switch o {
			case "candidate":
				channel = "candidate"
			case "edge":
				channel = "edge"
			case "beta":
				channel = "beta"
			case "classic":
				classic = true
			case "dangerous":
				dangerous = true
			case "devmode":
				devmode = true
			}
		}
		pkg, ok, err := datastores.Snap.Get(name)
		if err != nil {
			return nil, fmt.Sprintf("Could not check if %s is installed: %s", name, err), nil, StatusError, "", ""
		}
		if ok && pkg.Channel == channel && pkg.Classic == classic && pkg.Dangerous == dangerous && pkg.Devmode == devmode {
			return nil, fmt.Sprintf("%s is already installed", name), nil, StatusSuccess, "Version " + pkg.Version, ""
		}
		cmdArgs := []string{"install", name, "--" + channel}
		if classic {
			cmdArgs = append(cmdArgs, "--classic")
		}
		if dangerous {
			cmdArgs = append(cmdArgs, "--dangerous")
		}
		if devmode {
			cmdArgs = append(cmdArgs, "--devmode")
		}
		apply = func(out Output) bool {
			out.Runningf("Installing %s", name)
			if err := helper.Exec(nil, nil, "snap", cmdArgs...); err != nil {
				out.Errorf("Could not install %s: %s", name, helper.ExecErr(err))
				return false
			}
			out.Successf("Installed %s", name)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", name), apply, StatusInfo, "", ""
	},
	nil,
}
