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

type DpkgPackage struct {
	Name    string
	Version string
}

type dpkgPackagesList struct {
	initOnce sync.Once
	packages []DpkgPackage
}

func (d *dpkgPackagesList) Get(name string) (DpkgPackage, bool, error) {
	var err error
	d.initOnce.Do(func() { err = d.init() })
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	for _, pkg := range d.packages {
		if pkg.Name == name {
			return pkg, true, err
		}
	}
	return DpkgPackage{}, false, err
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
	pkg := DpkgPackage{}
	statusInstalled := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if statusInstalled {
				d.packages = append(d.packages, pkg)
			}
			statusInstalled = false
			pkg = DpkgPackage{}
		}
		info := strings.SplitN(line, ": ", 2)
		if len(info) != 2 {
			continue
		}
		switch info[0] {
		case "Package":
			pkg.Name = info[1]
		case "Version":
			pkg.Version = info[1]
		case "Status":
			if info[1] == "install ok installed" {
				statusInstalled = true
			}
		}
	}
	d.packages = append(d.packages, pkg)
	return scanner.Err()
}
