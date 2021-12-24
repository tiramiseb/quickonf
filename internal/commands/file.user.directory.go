package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tiramiseb/quickonf/internal/commands/shared"
)

func init() {
	register(fileUserDirectory)
}

var fileUserDirectory = Command{
	"file.user.directory",
	"Create a directory belonging to the given user (if path is relative, it is relative to the user's home directory",
	[]string{
		"Username",
		"Directory path",
	},
	nil,
	"Create Picz photos for alice\n  file.user.directory alice Picz",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		path := args[1]

		usr, err := shared.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join(usr.HomeDir, path)
		}

		info, err := os.Lstat(path)
		if err == nil {
			if info.IsDir() {
				return nil, fmt.Sprintf("Directory %s already exists", path), nil, StatusSuccess
			}
			return nil, fmt.Sprintf("%s already exists but is not a directory", path), nil, StatusError
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err.Error(), nil, StatusError
		}

		apply = &Apply{
			"file.user.directory",
			fmt.Sprintf("Will create directory %s", path),
			func(out Output) bool {
				if err := os.MkdirAll(path, 0755); err != nil {
					out.Errorf("Could not create directory %s: %s", path, err)
					return false
				}
				uid, err := strconv.Atoi(usr.Uid)
				if err != nil {
					out.Errorf("could not get UID: %s", err)
					return false
				}
				gid, err := strconv.Atoi(usr.Gid)
				if err != nil {
					out.Errorf("could not get GID: %s", err)
					return false
				}
				if err := os.Chown(path, uid, gid); err != nil {
					out.Errorf("Could not give ownership of %s to %s: %s", path, username, err)
					return false
				}
				out.Successf("Created directory %s", path)
				return true
			},
		}
		return nil, fmt.Sprintf("need to create directory %s", path), apply, StatusInfo
	},
	nil,
}
