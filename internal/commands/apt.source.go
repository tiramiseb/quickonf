package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

const aptSourcesBase = "/etc/apt/sources.list.d/"

func init() {
	register(aptSource)
}

var aptSource = Command{
	"apt.source",
	"Add apt source repository(ies), and update the available packages list if needed",
	[]string{
		"Local name of the source(s) (simple short text)",
		"Source(s)",
	},
	nil,
	"NextDNS\n  apt.source nextdns \"deb https://repo.nextdns.io/deb stable main\"",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]
		sources := args[1] + "\n"

		sourcesList := filepath.Join(aptSourcesBase, name+".list")
		existingB, err := os.ReadFile(sourcesList)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Sprintf("Could not read existing sources file: %s", err), nil, StatusError, "", ""
		}
		existing := string(existingB)
		if existing == sources {
			return nil, fmt.Sprintf("Sources %s already defined", name), nil, StatusSuccess, existing, sources
		}
		apply = func(out Output) bool {
			out.Runningf("Adding apt sources %s", name)
			if err := os.WriteFile(sourcesList, []byte(sources), 0o644); err != nil {
				out.Errorf("Could not write requested content to %s: %s", sourcesList, err)
				return false
			}
			out.Info("Waiting for dpkg to be available to update packages list")
			datastores.DpkgMutex.Lock()
			defer datastores.DpkgMutex.Unlock()
			out.Running("Updating packages list")
			if err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "update"); err != nil {
				out.Errorf("Could not update packages list: %s", helper.ExecErr(err))
				return false
			}
			out.Success("Updated packages list")
			return true
		}
		return nil, fmt.Sprintf("Need to add apt sources %s", name), apply, StatusInfo, existing, sources
	},
	nil,
}
