package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/commands/helper"
)

var Fish = fish{
	values: map[string]*fishUser{},
}

type fish struct {
	mutex  sync.Mutex
	values map[string]*fishUser
}

type fishUser struct {
	mutex            sync.Mutex
	initOnce         sync.Once
	user             User
	abbreviations    map[string]string
	variables        map[string]string
	pluginsInstalled []string
}

func (f *fish) Abbreviation(usr User, name string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	userData, ok := f.values[usr.Username]
	if !ok {
		userData = &fishUser{
			user:             usr,
			abbreviations:    map[string]string{},
			variables:        map[string]string{},
			pluginsInstalled: []string{},
		}
		f.values[usr.Username] = userData
	}
	return userData.abbreviation(name)
}

func (f *fish) Variable(usr User, name string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	userData, ok := f.values[usr.Username]
	if !ok {
		userData = &fishUser{
			user:             usr,
			abbreviations:    map[string]string{},
			variables:        map[string]string{},
			pluginsInstalled: []string{},
		}
		f.values[usr.Username] = userData
	}
	return userData.variable(name)
}

func (f *fish) IsPluginInstalled(usr User, name string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	userData, ok := f.values[usr.Username]
	if !ok {
		userData = &fishUser{
			user:             usr,
			abbreviations:    map[string]string{},
			variables:        map[string]string{},
			pluginsInstalled: []string{},
		}
		f.values[usr.Username] = userData
	}
	return userData.isPluginInstalled(name)
}

func (f *fish) Reset() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.values = map[string]*fishUser{}
}

func (f *fishUser) abbreviation(name string) (string, error) {
	var err error
	f.initOnce.Do(func() { err = f.init() })
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.abbreviations[name], err
}

func (f *fishUser) variable(name string) (string, error) {
	var err error
	f.initOnce.Do(func() { err = f.init() })
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.variables[name], err
}

func (f *fishUser) isPluginInstalled(name string) (bool, error) {
	var err error
	f.initOnce.Do(func() { err = f.init() })
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for _, p := range f.pluginsInstalled {
		if p == name {
			return true, err
		}
	}
	return false, err

}

func (f *fishUser) init() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Abbreviations
	var out bytes.Buffer
	if err := helper.ExecAs(f.user.User, nil, &out, "fish", "-c", "abbr --show"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		splat := strings.SplitN(line, " -- ", 2)
		if len(splat) != 2 {
			continue
		}
		abbr := strings.SplitN(splat[1], " ", 2)
		if len(abbr) != 2 {
			continue
		}
		if abbr[1][0] == '\'' && abbr[1][len(abbr[1])-1] == '\'' {
			f.abbreviations[abbr[0]] = abbr[1][1 : len(abbr[1])-1]
		} else {
			f.abbreviations[abbr[0]] = abbr[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Variables
	out.Reset()
	if err := helper.ExecAs(f.user.User, nil, &out, "fish", "-c", "set --long | grep -v '^history'"); err != nil {
		return err
	}
	scanner = bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		splat := strings.SplitN(line, " ", 2)
		if len(splat) != 2 {
			continue
		}
		f.variables[splat[0]] = splat[1]
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Plugins
	strBuilder := strings.Builder{}
	// Exit code 127 = command not found, therefore no plugin is installed
	if err := helper.ExecAs(f.user.User, nil, &strBuilder, "fish", "-c", "fisher list"); err != nil && helper.ExecErrCode(err) != 127 {
		return err
	}
	f.pluginsInstalled = strings.Split(strBuilder.String(), "\n")

	return nil
}
