package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands/datastores"
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
		lst, err := datastores.GitRemotes.List(uri)
		if err != nil {
			return nil, fmt.Sprintf("Could not list remote references for %s: %s", uri, err), nil, StatusError, "", ""
		}
		for _, l := range lst {
			if l.Name().Short() == ref {
				hash := l.Hash().String()
				return []string{hash}, fmt.Sprintf("Got hash for %s of %s", ref, uri), nil, StatusSuccess, "", hash
			}
		}
		return nil, fmt.Sprintf("Cannot find %s of %s", ref, uri), nil, StatusError, "", ""
	},
	datastores.GitRemotes.Reset,
}