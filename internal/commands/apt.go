package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

const dpkgStatusPath = "/var/lib/dpkg/status"

var (
	dpkgMutex    sync.Mutex
	dpkgPackages = dpkgPackagesList{}
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
		ok, err := dpkgPackages.installed(pkg)
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
				dpkgMutex.Lock()
				defer dpkgMutex.Unlock()
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
		dpkgPackages = dpkgPackagesList{}
	},
}

type dpkgPackage struct {
	name string
}

type dpkgPackagesList struct {
	initOnce sync.Once
	packages []dpkgPackage
}

func (d *dpkgPackagesList) installed(name string) (bool, error) {
	var err error
	d.initOnce.Do(func() { err = d.init() })
	for _, pkg := range d.packages {
		if pkg.name == name {
			return true, err
		}

	}
	return false, err
}

func (d *dpkgPackagesList) init() error {
	dpkgMutex.Lock()
	defer dpkgMutex.Unlock()
	f, err := os.Open(dpkgStatusPath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	pkg := dpkgPackage{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			d.packages = append(d.packages, pkg)
			pkg = dpkgPackage{}
		}
		info := strings.SplitN(line, ": ", 2)
		if len(info) != 2 {
			continue
		}
		switch info[0] {
		case "Package":
			pkg.name = info[1]
		}
	}
	return scanner.Err()
}
