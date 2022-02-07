package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(stackInstall)
}

var stackInstall = Command{
	"stack.install",
	"Execute the stack command (from Haskell)",
	[]string{
		"Directory where to execute the command",
	},
	nil,
	"Install foobar\n  stack.install /tmp/foobar",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]
		apply = &Apply{
			"stack.install",
			fmt.Sprintf("Will install from %s", path),
			func(out Output) bool {
				out.Runningf("Installing from %s", path)
				if err := helper.Exec(nil, nil, "stack", "--work-dir", path, "install"); err != nil {
					out.Errorf("Could not install from %s: %s", path, helper.ExecErr(err))
					return false
				}
				out.Successf("Installed from %s", path)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to install from %s", path), apply, StatusInfo
	},
	nil,
}
