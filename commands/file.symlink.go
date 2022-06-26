package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	register(fileSymlink)
}

var fileSymlink = &Command{
	"file.symlink",
	"Create a symlink",
	[]string{
		"Absolute path of the symlink",
		"Absolute path of the target file",
	},
	nil,
	"Very dumb symlink\n  file.symlink /home/alice/temp /tmp",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		link := args[0]
		target := args[1]
		if !filepath.IsAbs(link) {
			return nil, fmt.Sprintf("%s is not an absolute path", link), nil, StatusError, "", ""
		}
		if !filepath.IsAbs(target) {
			return nil, fmt.Sprintf("%s is not an absolute path", target), nil, StatusError, "", ""
		}

		var (
			needMessage   string
			mustBeRemoved bool
		)
		info, err := os.Lstat(link)
		switch {
		case err != nil:
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err.Error(), nil, StatusError, "", ""
			}
			needMessage = fmt.Sprintf("Need to create %s", link)
		case info.Mode()&os.ModeDir != 0:
			needMessage = fmt.Sprintf("Need to remove directory %s and create link", link)
			mustBeRemoved = true
		case info.Mode()&os.ModeSymlink == 0:
			needMessage = fmt.Sprintf("Need to remove file %s and create link", link)
			mustBeRemoved = true
		default:
			existingTarget, err := filepath.EvalSymlinks(link)
			if err != nil {
				return nil, err.Error(), nil, StatusError, "", ""
			}
			if target == existingTarget {
				return nil, fmt.Sprintf("%s already targets %s", link, target), nil, StatusSuccess, "", ""
			}
			needMessage = fmt.Sprintf("Need to remove link %s and recreate it with target %s", link, target)
			mustBeRemoved = true
		}

		apply = func(out Output) bool {
			if mustBeRemoved {
				out.Runningf("Removing %s", link)
				if err := os.RemoveAll(link); err != nil {
					out.Errorf("Could not remove %s: %s", link, err)
					return false
				}
			}
			out.Runningf("Creating link %s to %s", link, target)
			if err := os.Symlink(target, link); err != nil {
				out.Errorf("Could not create %s: %s", link, err)
				return false
			}
			out.Successf("Link %s created with target %s", link, target)
			return true
		}

		return nil, needMessage, apply, StatusInfo, "", ""
	},
	nil,
}
