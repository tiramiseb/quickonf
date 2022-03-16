package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
)

func init() {
	register(userFileDirectory)
}

var userFileDirectory = Command{
	"user.file.directory",
	"Create a directory belonging to the given user (if path is relative, it is relative to the user's home directory",
	[]string{
		"Username",
		"Directory path",
	},
	nil,
	"Create Picz photos for alice\n  user.file.directory alice Picz",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		path := args[1]

		usr, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join(usr.User.HomeDir, path)
		}

		info, err := os.Lstat(path)
		if err == nil {
			if info.IsDir() {
				return nil, fmt.Sprintf("Directory %s already exists", path), nil, StatusSuccess, "", ""
			}
			return nil, fmt.Sprintf("%s already exists but is not a directory", path), nil, StatusError, "", ""
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err.Error(), nil, StatusError, "", ""
		}

		apply = func(out Output) bool {
			out.Runningf("Creating directory %s", path)
			if err := os.MkdirAll(path, 0o755); err != nil {
				out.Errorf("Could not create directory %s: %s", path, err)
				return false
			}
			if err := os.Chown(path, usr.Uid, usr.Group.Gid); err != nil {
				out.Errorf("Could not give ownership of %s to %s: %s", path, username, err)
				return false
			}
			out.Successf("Created directory %s", path)
			return true
		}
		return nil, fmt.Sprintf("Need to create directory %s", path), apply, StatusInfo, "", ""
	},
	datastores.Users.Reset,
}
