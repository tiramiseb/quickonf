package datastores

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
)

const groupsfile = "/etc/group"

type User struct {
	User *user.User

	Username string
	Uid      int

	// Main group
	Group Group

	// Other groups
	Groups []Group
}

type Group struct {
	Name string
	Gid  int

	Members []string
}

type users struct {
	usermutex  sync.Mutex
	groupmutex sync.Mutex
	initOnce   sync.Once
	users      map[string]User
	groups     []Group
}

// Base of users and groups on the system
var (
	FakeUserForSystem = User{Username: "|system|"}
	Users             = &users{users: map[string]User{}}
)

func (u *users) Get(username string) (usr User, err error) {
	u.usermutex.Lock()
	defer u.usermutex.Unlock()
	usr, ok := u.users[username]
	if !ok {
		usr.Username = username

		lookedupUser, err := user.Lookup(username)
		if err != nil {
			return usr, err
		}
		usr.User = lookedupUser

		uid, err := strconv.Atoi(lookedupUser.Uid)
		if err != nil {
			return usr, err
		}
		usr.Uid = uid

		gid, err := strconv.Atoi(lookedupUser.Gid)
		if err != nil {
			return usr, err
		}
		group, ok, err := u.GetGroupByID(gid)
		if err != nil {
			return usr, err
		}
		if !ok {
			return usr, fmt.Errorf("Group with ID %d does not exist", gid)
		}
		usr.Group = group

		for _, g := range u.groups {
			if g.HasUser(username) {
				usr.Groups = append(usr.Groups, g)
			}
		}

		u.users[username] = usr
	}
	return usr, nil
}

func (u *users) GetGroup(name string) (group Group, ok bool, err error) {
	u.initOnce.Do(func() { err = u.initGroups() })
	u.groupmutex.Lock()
	defer u.groupmutex.Unlock()
	for _, g := range u.groups {
		if g.Name == name {
			group = g
			ok = true
		}
	}
	return

}

func (u *users) GetGroupByID(gid int) (group Group, ok bool, err error) {
	u.initOnce.Do(func() { err = u.initGroups() })
	u.groupmutex.Lock()
	defer u.groupmutex.Unlock()
	for _, g := range u.groups {
		if g.Gid == gid {
			group = g
			ok = true
		}
	}
	return

}

func (u *users) Reset() {
	u.usermutex.Lock()
	u.groupmutex.Lock()
	defer u.usermutex.Unlock()
	defer u.groupmutex.Unlock()
	u.users = map[string]User{}
	u.groups = u.groups[:0]
	u.initOnce = sync.Once{}
}

func (u *users) initGroups() error {
	u.groupmutex.Lock()
	defer u.groupmutex.Unlock()
	f, err := os.Open(groupsfile)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) != 4 {
			continue
		}
		gid, err := strconv.Atoi(line[2])
		if err != nil {
			return err
		}
		u.groups = append(u.groups, Group{
			Name:    line[0],
			Gid:     gid,
			Members: strings.Split(line[3], ","),
		})
	}
	return scanner.Err()
}

func (g Group) HasUser(username string) bool {
	for _, u := range g.Members {
		if u == username {
			return true
		}
	}
	return false
}
