package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(makefile)
}

var makefile = &Command{
	"make",
	"Execute the make command",
	[]string{
		"Directory where to execute the command",
		"Make target",
	},
	nil,
	"Compile foobar\n  make /tmp/foobar compile",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		path := args[0]
		target := args[1]
		apply = func(out Output) bool {
			out.Runningf("Making %s in %s", target, path)
			if err := helper.Exec(nil, nil, "make", "-C", path, target); err != nil {
				out.Errorf("Could not make %s in %s: %s", target, path, helper.ExecErr(err))
				return false
			}
			out.Successf("Made %s in %s", target, path)
			return true
		}
		return nil, fmt.Sprintf("Need to make %s in %s", target, path), apply, StatusInfo, "", ""
	},
	nil,
}
