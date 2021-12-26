package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	register(fileAbsent)
}

var fileAbsent = Command{
	"file.absent",
	"Make sure a file is absent",
	[]string{
		"Absolute path of the file to remove (including directory with content)",
	},
	nil,
	`Make sure there is no "Photos" directory for jack\n  file.absent /home/jack/Photos`,
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError
		}
		_, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("%s does not exist", path), nil, StatusSuccess
			}
			return nil, err.Error(), nil, StatusError
		}

		apply = &Apply{
			"file.absent",
			fmt.Sprintf("Will remove %s", path),
			func(out Output) bool {
				if err := os.RemoveAll(path); err != nil {
					out.Errorf("could not remove %s: %s", path, err)
					return false
				}
				out.Successf("Removed %s", path)
				return true
			},
		}

		return nil, fmt.Sprintf("need to remove %s", path), apply, StatusInfo

	},
	nil,
}
