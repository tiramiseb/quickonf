package commands

import (
	"fmt"
	"strings"

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
		"Make targets...",
	},
	nil,
	"Compile foobar\n  make /tmp/foobar compile",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		path := args[0]
		targets := args[1:]
		apply = func(out Output) bool {
			out.Runningf("Making %s in %s", strings.Join(targets, " "), path)
			args = append([]string{"-C", path}, targets...)
			if err := helper.Exec(nil, nil, "make", args...); err != nil {
				out.Errorf("Could not make %s in %s: %s", strings.Join(targets, " "), path, helper.ExecErr(err))
				return false
			}
			out.Successf("Made %s in %s", strings.Join(targets, " "), path)
			return true
		}
		return nil, fmt.Sprintf("Need to make %s in %s", strings.Join(targets, " "), path), apply, StatusInfo, "", ""
	},
	nil,
}
