package datastores

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

var (
	GoEnv = goenv{
		values: map[string]*goenvUser{},
	}
)

type goenvUser struct {
	mutex    sync.Mutex
	initOnce sync.Once
	user     User
	values   map[string]string
}

type goenv struct {
	mutex  sync.Mutex
	values map[string]*goenvUser
}

func (g *goenv) Get(usr User, variable string) (string, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	userData, ok := g.values[usr.Username]
	if !ok {
		userData = &goenvUser{
			user:   usr,
			values: map[string]string{},
		}
		g.values[usr.Username] = userData
	}
	return userData.get(variable)
}

func (g *goenv) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.values = map[string]*goenvUser{}
}

func (g *goenvUser) get(variable string) (string, error) {
	var err error
	g.initOnce.Do(func() { err = g.init() })
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.values[variable], err
}

func (g *goenvUser) init() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.ExecAs(g.user.User, nil, &out, "go", "env", "-json"); err != nil {
		return err
	}
	dec := json.NewDecoder(&out)
	return dec.Decode(&(g.values))
}
