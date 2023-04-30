package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(cmdOut)
}

var cmdOut = &Command{
	"cmd.out",
	"Get the output of a command (use sparingly - only for simple commands when a Quickonf command does not exist) - if command does not exist, output is the empty string",
	[]string{
		"Command name",
		"Command arguments",
	},
	[]string{
		"Output of the command",
	},
	"Install Go\n  out = cmd.out go version\n  version = regexp.submatch \"go version go([0-9]+\\.[0-9]+\\.[0-9]+)\" <out>\n  [...]",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		cmd := args[0]
		arguments := args[1]
		var buf bytes.Buffer
		if err := helper.Exec(nil, &buf, cmd, arguments); err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				return []string{""}, fmt.Sprintf("%s does not exist", cmd), nil, StatusSuccess, "", ""
			}
			return []string{""}, fmt.Sprintf("Could not execute %s: %s", cmd, helper.ExecErr(err)), nil, StatusError, "", ""
		}
		out := buf.String()
		return []string{out}, fmt.Sprintf("Executed %s", cmd), nil, StatusSuccess, "", out
	},
	nil,
}
