package datastores

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/tiramiseb/quickonf/commands/helper"
)

// List of alternatives
var Alternatives = alternatives{alternatives: map[string]string{}}

type alternatives struct {
	mutex        sync.Mutex
	initOnce     sync.Once
	alternatives map[string]string
}

func (a *alternatives) Get(name string) (string, error) {
	var err error
	a.initOnce.Do(func() { err = a.init() })
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.alternatives[name], err
}

func (a *alternatives) Reset() {
	a.alternatives = map[string]string{}
	a.initOnce = sync.Once{}
}

func (a *alternatives) init() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var out bytes.Buffer
	if err := helper.Exec(nil, &out, "update-alternatives", "--get-selections"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		if len(data) != 3 {
			continue
		}
		a.alternatives[data[0]] = data[2]
	}
	return scanner.Err()
}
