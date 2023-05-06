package datastores

import (
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/commands/helper"
)

var VSCodeExtensions = vscodeExtensions{
	values: map[string]*userVSCodeExtensions{},
}

type userVSCodeExtensions struct {
	mutex               sync.Mutex
	initOnce            sync.Once
	user                User
	installedExtensions []string
}

type vscodeExtensions struct {
	mutex  sync.Mutex
	values map[string]*userVSCodeExtensions
}

func (v *vscodeExtensions) Installed(usr User, id string) (bool, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	userData, ok := v.values[usr.Username]
	if !ok {
		userData = &userVSCodeExtensions{
			user: usr,
		}
		v.values[usr.Username] = userData
	}
	return userData.installed(id)
}

func (v *vscodeExtensions) Reset() {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.values = map[string]*userVSCodeExtensions{}
}

func (u *userVSCodeExtensions) installed(id string) (bool, error) {
	var err error
	u.initOnce.Do(func() { err = u.init() })
	u.mutex.Lock()
	defer u.mutex.Unlock()
	id = strings.ToLower(id)
	for _, thisID := range u.installedExtensions {
		if id == thisID {
			return true, err
		}
	}
	return false, err
}

func (u *userVSCodeExtensions) init() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.ExecAs(u.user.User, nil, &out, "code", "--list-extensions"); err == nil {
		u.installedExtensions = strings.Split(strings.ToLower(out.String()), "\n")
	}
	return nil
}
