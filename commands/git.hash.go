package commands

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(gitHash)
}

var gitHash = Command{
	"git.hash",
	"Get the latest hash for a reference (generally branch or tag) in a git repository",
	[]string{
		"Git repository URI",
		"Reference (generally, branch or tag)",
	},
	[]string{
		"Commit hash",
	},
	"Pop Shell\n  pophash = git.hash https://github.com/pop-os/shell.git master",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uri := args[0]
		ref := args[1]
		var out strings.Builder
		if err := helper.Exec(nil, &out, "git", "ls-remote", uri, ref); err != nil {
			return nil, fmt.Sprintf("Could not list references for %s: %s", uri, helper.ExecErr(err)), nil, StatusError, "", ""
		}
		fields := strings.Fields(out.String())
		if len(fields) != 2 {
			return nil, fmt.Sprintf("Wrong format for references %s: %s", uri, out.String()), nil, StatusError, "", ""
		}
		return []string{fields[0]}, fmt.Sprintf("Got hash for %s of %s", ref, uri), nil, StatusSuccess, "", fields[0]
	},
	nil,
}
