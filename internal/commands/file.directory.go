package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	register(fileDirectory)
}

var fileDirectory = Command{
	"file.directory",
	"Create a directory (path must be absolute)",
	[]string{
		"Directory path",
	},
	nil,
	"Create web root\n  file.directory /srv/web",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]

		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError
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
			"file.directory",
			fmt.Sprintf("Will create directory %s", path),
			func(out Output) bool {
				out.Runningf("Creating directory %s", path)
				if err := os.MkdirAll(path, 0o755); err != nil {
					out.Errorf("Could not create directory %s: %s", path, err)
					return false
				}
				out.Successf("Created directory %s", path)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to create directory %s", path), apply, StatusInfo
	},
	nil,
}
