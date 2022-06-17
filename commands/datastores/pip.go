package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/commands/helper"
)

var PipPackages = pip{
	packages: map[string]PipPackage{},
}

type PipPackage struct {
	Name    string
	Version string
}

type pip struct {
	mutex    sync.Mutex
	initOnce sync.Once
	packages map[string]PipPackage
}

func (p *pip) Get(name string) (PipPackage, bool, error) {
	var err error
	p.initOnce.Do(func() { err = p.init() })
	p.mutex.Lock()
	defer p.mutex.Unlock()
	pkg, ok := p.packages[name]
	return pkg, ok, err
}

func (p *pip) Reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.packages = map[string]PipPackage{}
	p.initOnce = sync.Once{}
}

func (p *pip) init() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.Exec(nil, &out, "pip3", "list"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) != 2 {
			continue
		}
		if line[0] == "Package" || line[0][0] == '-' {
			continue
		}
		pkg := PipPackage{
			Name:    line[0],
			Version: line[1],
		}
		p.packages[line[0]] = pkg

	}
	return scanner.Err()
}
