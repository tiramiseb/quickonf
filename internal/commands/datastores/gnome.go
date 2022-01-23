package datastores

import (
	"encoding/json"
	"strings"
	"sync"
)

var GnomeExtensions = gnomeExtensions{
	values: map[string]*userGnomeExtensions{},
}

type userGnomeExtensions struct {
	mutex             sync.Mutex
	initOnce          sync.Once
	user              User
	enabledExtensions []string
}

type gnomeExtensions struct {
	mutex    sync.Mutex
	initOnce sync.Once
	values   map[string]*userGnomeExtensions
}

func (g *gnomeExtensions) Enabled(usr User, uuid string) (bool, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	userData, ok := g.values[usr.Username]
	if !ok {
		userData = &userGnomeExtensions{
			user: usr,
		}
		g.values[usr.Username] = userData
	}
	return userData.enabled(uuid)
}

func (g *gnomeExtensions) Reset() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.values = map[string]*userGnomeExtensions{}
}

func (u *userGnomeExtensions) enabled(uuid string) (bool, error) {
	var err error
	u.initOnce.Do(func() { err = u.init() })
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, ext := range u.enabledExtensions {
		if ext == uuid {
			return true, err
		}
	}
	return false, err
}

func (u *userGnomeExtensions) init() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	extensionsStr, err := Dconf.Get(u.user, "/org/gnome/shell/enabled-extensions")
	if err != nil {
		return err
	}
	extensionsStr = strings.ReplaceAll(extensionsStr, "'", "\"")
	return json.Unmarshal([]byte(extensionsStr), &u.enabledExtensions)
}
