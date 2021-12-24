package commands

import (
	"fmt"
	"os"
	"path/filepath"

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
		}
		if !os.IsNotExist(err) {
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
				out.Successf("Created directory %s", path)
				return true
			},
		}
		return nil, fmt.Sprintf("need to create directory %s", path), apply, StatusInfo
	},
	nil,
}
