package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

var GitConfig = gitConfigs{
	values: map[string]*gitConfig{},
}

type gitConfig struct {
	user     User
	mutex    sync.Mutex
	initOnce sync.Once
	values   map[string]string
}

type gitConfigs struct {
	mutex  sync.Mutex
	values map[string]*gitConfig // Map of username to list of values
}

func (g *gitConfigs) Get(user User, key string) (string, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	userData, ok := g.values[user.Username]
	if !ok {
		userData = &gitConfig{
			user:   user,
			values: map[string]string{},
		}
		g.values[user.Username] = userData
	}
	return userData.get(key)
}

func (g *gitConfigs) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.values = map[string]*gitConfig{}
}

func (g *gitConfig) get(name string) (string, error) {
	var err error
	g.initOnce.Do(func() { err = g.init() })
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.values[name], err
}

func (g *gitConfig) init() (err error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var out bytes.Buffer
	if g.user.Username == FakeUserForSystem.Username {
		if err := helper.Exec(nil, &out, "git", "config", "--system", "--list"); err != nil {
			return err
		}
	} else {
		if err := helper.ExecAs(g.user.User, nil, &out, "git", "config", "--global", "--list"); err != nil {
			return err
		}
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) == 2 {
			g.values[line[0]] = line[1]
		}
	}
	return scanner.Err()
}
