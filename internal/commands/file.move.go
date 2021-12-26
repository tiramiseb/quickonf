package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	register(fileDirectoryMove)
}

var fileDirectoryMove = Command{
	"file.move",
	"Move a file or directory to another place (does not fail if the source does not exist)",
	[]string{
		"Souce path",
		"Destination path",
	},
	nil,
	"Move previous documents directory for jane\n  file.move /home/jane.old/Documents /home/jane/Documents",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		source := args[0]
		destination := args[1]
		if !filepath.IsAbs(source) {
			return nil, fmt.Sprintf("%s is not an absolute path", source), nil, StatusError
		}
		if !filepath.IsAbs(destination) {
			return nil, fmt.Sprintf("%s is not an absolute path", destination), nil, StatusError
		}
		if _, err := os.Lstat(source); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("Directory %s does not exist", source), nil, StatusSuccess
			}
			return nil, err.Error(), nil, StatusError
		}
		_, err := os.Lstat(destination)
		if err == nil {
			return nil, fmt.Sprintf("Directory %s already exists", destination), nil, StatusError
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err.Error(), nil, StatusError
		}

		apply = &Apply{
			"file.move",
			fmt.Sprintf("Will move %s to %s", source, destination),
			func(out Output) bool {
				if err := os.Rename(source, destination); err != nil {
					out.Errorf("Could not move %s to %s: %s", source, destination, err)
					return false
				}
				out.Successf("Moved %s to %s", source, destination)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to move %s to %s", source, destination), apply, StatusInfo
	},
	nil,
}
