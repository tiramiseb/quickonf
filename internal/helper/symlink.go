package helper

import (
	"fmt"
	"os"
)

// SymlinkStatus is the status after trying to create a symbolic link
type SymlinkStatus int

// Symlink status is returned when trying to create a symbolic link
const (
	SymlinkError SymlinkStatus = iota
	SymlinkAleradyExists
	SymlinkCreated
)

// Symlink creates a symbolic link
func Symlink(path, target string) (status SymlinkStatus, err error) {
	if stat, err := os.Lstat(path); err == nil {
		if stat.Mode()&os.ModeSymlink != os.ModeSymlink {
			return SymlinkError, fmt.Errorf("%s already exists but is not a symlink", path)
		}
		if currentTarget, err := os.Readlink(path); err == nil {
			if currentTarget == target {
				return SymlinkAleradyExists, nil
			}
			if !Dryrun {
				if err := os.Remove(path); err != nil {
					return SymlinkError, err
				}
			}
		} else if !os.IsNotExist(err) {
			return SymlinkError, err
		}

	} else if !os.IsNotExist(err) {
		return SymlinkError, err
	}
	if !Dryrun {
		if err := os.Symlink(target, path); err != nil {
			return SymlinkError, err
		}
	}
	return SymlinkCreated, nil
}
