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

var fileAbsent = &Command{
	"file.absent",
	"Make sure a file is absent",
	[]string{
		"Absolute path of the file to remove (including directory with content)",
	},
	nil,
	"Make sure there is no \"Photos\" directory for jack\n  file.absent /home/jack/Photos",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		path := args[0]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError, "", ""
		}
		_, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("%s does not exist", path), nil, StatusSuccess, "", ""
			}
			return nil, err.Error(), nil, StatusError, "", ""
		}

		apply = func(out Output) bool {
			out.Runningf("Removing %s", path)
			if err := os.RemoveAll(path); err != nil {
				out.Errorf("Could not remove %s: %s", path, err)
				return false
			}
			out.Successf("Removed %s", path)
			return true
		}

		return nil, fmt.Sprintf("Need to remove %s", path), apply, StatusInfo, "", ""
	},
	nil,
}
