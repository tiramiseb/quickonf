package helper

import (
	"fmt"
	"os"
)

// Symlink creates a symbolic link.
//
// If the symlink already exists and points to the specified target, return ResultAlready.
//
// If in dry-run mode, do not create nor change the symlink.
func Symlink(path, target string) (status ResultStatus, err error) {
	if stat, err := os.Lstat(path); err == nil {
		if stat.Mode()&os.ModeSymlink != os.ModeSymlink {
			return ResultError, fmt.Errorf("%s already exists but is not a symlink", path)
		}
		if currentTarget, err := os.Readlink(path); err == nil {
			if currentTarget == target {
				return ResultAlready, nil
			}
			if !Dryrun {
				if err := os.Remove(path); err != nil {
					return ResultError, err
				}
			}
		} else if !os.IsNotExist(err) {
			return ResultError, err
		}

	} else if !os.IsNotExist(err) {
		return ResultError, err
	}
	if Dryrun {
		return ResultDryrun, nil
	}
	if err := os.Symlink(target, path); err != nil {
		return ResultError, err
	}
	return ResultSuccess, nil
}
