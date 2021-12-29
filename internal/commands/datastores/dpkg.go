package datastores

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const dpkgStatusPath = "/var/lib/dpkg/status"

var (
	// Lock this mutex whenever you start manipulating apt/dpkg packages and stuff
	DpkgMutex sync.Mutex
	// List of installed packages
	DpkgPackages = dpkgPackagesList{}
)

type dpkgPackage struct {
	name string
}

type dpkgPackagesList struct {
	initOnce sync.Once
	packages []dpkgPackage
}

func (d *dpkgPackagesList) Installed(name string) (bool, error) {
	var err error
	d.initOnce.Do(func() { err = d.init() })
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	for _, pkg := range d.packages {
		if pkg.name == name {
			return true, err
		}
	}
	return false, err
}

func (d *dpkgPackagesList) Reset() {
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	d.packages = nil
	d.initOnce = sync.Once{}
}

func (d *dpkgPackagesList) init() error {
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	f, err := os.Open(dpkgStatusPath)
	if err != nil {
		return err
	}
	defer f.Close()
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
	d.packages = append(d.packages, pkg)
	return scanner.Err()
}
