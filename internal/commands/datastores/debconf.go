package datastores

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const debconfDB = "/var/cache/debconf/config.dat"

var (
	Debconf = debconf{}
)

type DebconfParameter struct {
	Name   string
	Value  string
	Owners []string
}

type debconf struct {
	initOnce   sync.Once
	parameters []DebconfParameter
}

func (d *debconf) Get(pkg, name string) (param DebconfParameter, exists bool, err error) {
	d.initOnce.Do(func() { err = d.init() })
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
paramsLoop:
	for _, p := range d.parameters {
		if p.Name == name {
			for _, owner := range p.Owners {
				if owner == pkg {
					param = p
					exists = true
					break paramsLoop
				}
			}
		}
	}
	return
}

func (d *debconf) Reset() {
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	d.parameters = nil
	d.initOnce = sync.Once{}
}

func (d *debconf) init() error {
	DpkgMutex.Lock()
	defer DpkgMutex.Unlock()
	f, err := os.Open(debconfDB)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	param := DebconfParameter{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			d.parameters = append(d.parameters, param)
			param = DebconfParameter{}
		}
		info := strings.SplitN(line, ": ", 2)
		if len(info) != 2 {
			continue
		}
		switch info[0] {
		case "Name":
			param.Name = info[1]
		case "Value":
			param.Value = info[1]
		case "Owners":
			param.Owners = strings.Split(info[1], ", ")
		}
	}
	d.parameters = append(d.parameters, param)
	return scanner.Err()
}
