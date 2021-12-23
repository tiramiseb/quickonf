package commands

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

var (
	aptMutex    sync.Mutex
	aptPackages = aptPackagesList{}
)

func init() {
	register(apt)
}

var apt = Command{
	"apt",
	"Install a package using apt",
	[]string{"Name of the package to install"},
	nil,
	"Install the \"ipcalc\" tool\n  apt ipcalc",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		pkg := args[0]
		ok, err := aptPackages.installed(pkg)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		if ok {
			return nil, fmt.Sprintf("%s is already installed", pkg), nil, StatusSuccess
		}

		apply = &Apply{
			"apt",
			fmt.Sprintf("will install %s", pkg),
			func(out Output) bool {
				out.Infof("waiting for apt to be available to install %s", pkg)
				aptMutex.Lock()
				defer aptMutex.Unlock()
				wait, err := helper.Exec([]string{"DEBIAN_FRONTEND=noninteractive"}, nil, "apt-get", "--yes", "--quiet", "install", pkg)
				if err != nil {
					out.Errorf("could not install %s: %s", pkg, err)
					return false
				}
				out.Infof("installing %s", pkg)
				if err := wait(); err != nil {
					out.Errorf("could not install %s: %s", pkg, err)
					return false
				}
				out.Successf("installed %s", pkg)
				return true
			},
		}
		return nil, fmt.Sprintf("need to install %s", pkg), apply, StatusInfo
	},
	func() {
		aptPackages = aptPackagesList{}
	},
}

type aptPackagesList struct {
	initOnce sync.Once
	content  *helper.SearchableStrings
}

func (a *aptPackagesList) installed(pkg string) (bool, error) {
	var err error
	a.initOnce.Do(func() { err = a.init() })
	return a.content.Contains(pkg), err
}

func (a *aptPackagesList) init() error {
	a.content = &helper.SearchableStrings{}
	lines, err := helper.ExecOutAsLines(nil, "dpkg", "--get-selections")
	if err != nil {
		return fmt.Errorf("could not get list of installed packages: %w", err)
	}
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) != 2 {
			continue
		}
		if fields[1] == "install" {
			a.content.Add(fields[0])
		}
	}
	return nil
}
