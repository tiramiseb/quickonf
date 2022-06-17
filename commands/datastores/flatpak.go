package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/commands/helper"
)

var Flatpak = flatpak{}

type FlatpakPackage struct {
	ApplicationID string
}

type flatpak struct {
	mutex    sync.Mutex
	initOnce sync.Once
	packages []FlatpakPackage
}

func (f *flatpak) Get(applicationID string) (FlatpakPackage, bool, error) {
	var err error
	f.initOnce.Do(func() { err = f.init() })
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for _, pkg := range f.packages {
		if pkg.ApplicationID == applicationID {
			return pkg, true, err
		}
	}
	return FlatpakPackage{}, false, err
}

func (f *flatpak) Reset() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.packages = []FlatpakPackage{}
	f.initOnce = sync.Once{}
}

func (f *flatpak) init() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.Exec(nil, &out, "flatpak", "list", "--columns=application"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) != 1 {
			continue
		}
		if line[0] == "Application" {
			continue
		}
		pkg := FlatpakPackage{
			ApplicationID: line[0],
		}
		f.packages = append(f.packages, pkg)
	}
	return nil
}
