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
	`NextDNS\n  apt.source nextdns "deb https://repo.nextdns.io/deb stable main"`,
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		name := args[0]
		sources := args[1] + "\n"

		sourcesList := filepath.Join(aptSourcesBase, name+".list")
		existingB, err := os.ReadFile(sourcesList)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("Could not read existing sources file: %s", err), nil, StatusError
			}
			err = nil
		}
		existing := string(existingB)
		if existing == sources {
			return nil, fmt.Sprintf("Sources %s already defined", name), nil, StatusSuccess
		}
		apply = &Apply{
			"apt.source",
			fmt.Sprintf("Will add apt sources %s", name),
			func(out Output) bool {
				out.Infof("Adding apt sources %s", name)
				if err := os.WriteFile(sourcesList, []byte(sources), 0644); err != nil {
					out.Errorf("Could not write requested content to %s: %s", sourcesList, err)
					return false
				}
				out.Info("Waiting for dpkg to be available to update packages list")
				datastores.DpkgMutex.Lock()
				defer datastores.DpkgMutex.Unlock()
				out.Infof("Updating packages list")
				wait, err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "update")
				if err != nil {
					out.Errorf("Could not update packages list: %s", err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not update packages list: %s", err)
					return false
				}
				return true
			},
		}
		return nil, fmt.Sprintf("Need to add apt sources %s", name), apply, StatusInfo
	},
	nil,
}