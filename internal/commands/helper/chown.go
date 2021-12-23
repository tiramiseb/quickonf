package helper

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

func Chown(file string, usr *user.User) error {
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return fmt.Errorf("could not get user ID: %v", err)
	}
	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return fmt.Errorf("could not get group ID: %v", err)
	}
	if err := os.Chown(file, uid, gid); err != nil {
		return fmt.Errorf("could not change ownership of %s: %v", file, err)
	}
	return nil
}
