package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

var Snap = snap{
	packages: map[string]SnapPackage{},
}

type SnapPackage struct {
	Name      string
	Channel   string
	Classic   bool
	Devmode   bool
	Dangerous bool
}

type snap struct {
	mutex    sync.Mutex
	initOnce sync.Once
	packages map[string]SnapPackage
}

func (s *snap) Get(name string) (SnapPackage, bool, error) {
	var err error
	s.initOnce.Do(func() { err = s.init() })
	s.mutex.Lock()
	defer s.mutex.Unlock()
	pkg, ok := s.packages[name]
	return pkg, ok, err
}

func (s *snap) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.packages = map[string]SnapPackage{}
	s.initOnce = sync.Once{}
}

func (s *snap) init() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.Exec(nil, &out, "snap", "list"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) != 6 {
			continue
		}
		if line[0] == "Name" {
			continue
		}
		pkg := SnapPackage{
			Name: line[0],
		}
		tracking := strings.FieldsFunc(line[3], func(r rune) bool {
			return r == '/'
		})
		if len(tracking) >= 2 {
			pkg.Channel = tracking[1]
		}
		switch line[5] {
		case "classic":
			pkg.Classic = true
		case "dangerous":
			pkg.Dangerous = true
		case "devmode":
			pkg.Devmode = true
		}
		s.packages[line[0]] = pkg

	}
	return nil
}
