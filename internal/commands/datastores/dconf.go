package datastores

import (
	"bytes"
	"sync"

	"gopkg.in/ini.v1"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

var (
	Dconf = dconf{
		values: map[string]*dconfUser{},
	}
)

type dconfUser struct {
	mutex    sync.Mutex
	initOnce sync.Once
	user     User
	values   map[string]string
}

type dconf struct {
	mutex  sync.Mutex
	values map[string]*dconfUser
}

func (d *dconf) Get(usr User, key string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	userData, ok := d.values[usr.Username]
	if !ok {
		userData = &dconfUser{
			user:   usr,
			values: map[string]string{},
		}
		d.values[usr.Username] = userData
	}
	return userData.get(key)
}

func (d *dconf) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.values = map[string]*dconfUser{}
}

func (d *dconfUser) get(key string) (string, error) {
	var err error
	d.initOnce.Do(func() { err = d.init() })
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.values[key], err
}

func (d *dconfUser) init() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.ExecAs(d.user.User, nil, &out, "dconf", "dump", "/"); err != nil {
		return err
	}
	f, err := ini.Load(&out)
	if err != nil {
		return err
	}
	for _, section := range f.Sections() {
		prefix := section.Name()
		for key, value := range section.KeysHash() {
			d.values["/"+prefix+"/"+key] = value
		}
	}
	return nil
}
