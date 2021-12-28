package datastores

import "os/user"

type users map[string]*user.User

// Base of users on the system
var Users = users{}

func (u users) Get(username string) (usr *user.User, err error) {
	usr, ok := u[username]
	if !ok {
		usr, err = user.Lookup(username)
		u[username] = usr
	}
	return
}
