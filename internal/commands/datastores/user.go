package datastores

import (
	"os/user"
	"sync"
)

type users map[string]*user.User

// Base of users on the system
var Users = users{}
var usersMutex sync.Mutex

func (u users) Get(username string) (usr *user.User, err error) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	usr, ok := u[username]
	if !ok {
		usr, err = user.Lookup(username)
		u[username] = usr
	}
	return
}
