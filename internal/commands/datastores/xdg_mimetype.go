package datastores

import (
	"sync"

	"gopkg.in/ini.v1"
)

const xdgMimetypeFile = "/etc/xdg/mimeapps.list"

var (
	XdgMimetypes = xdgMimetypes{}
)

type xdgMimetypes struct {
	mutex    sync.Mutex
	initOnce sync.Once
	defaults map[string]string
}

func (x *xdgMimetypes) Get(mimetype string) (string, error) {
	var err error
	x.initOnce.Do(func() { err = x.init() })
	x.mutex.Lock()
	defer x.mutex.Unlock()
	return x.defaults[mimetype], err
}

func (x *xdgMimetypes) Reset() {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	x.defaults = map[string]string{}
	x.initOnce = sync.Once{}
}

func (x *xdgMimetypes) init() error {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	x.defaults = map[string]string{}

	conf, err := ini.LooseLoad(xdgMimetypeFile)
	if err != nil {
		return err
	}
	section := conf.Section("Default Applications")
	for _, k := range section.Keys() {
		x.defaults[k.Name()] = k.String()
	}
	return nil
}
