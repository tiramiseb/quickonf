package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	register(fileDirectoryMoveNocheck)
}

var fileDirectoryMoveNocheck = Command{
	"file.move.nocheck",
	"Move a file or directory to another place (with no check). Necessary if the file does not exist in the check phase (for instance if it is downloaded when applying).",
	[]string{
		"Souce path",
		"Destination path",
	},
	nil,
	"Move previous documents directory for jane\n  file.move.nocheck /home/jane.old/Documents /home/jane/Documents",
	func(args []string) (result []string, msg string, apply Apply, status Status) {
		source := args[0]
		destination := args[1]
		if !filepath.IsAbs(source) {
			return nil, fmt.Sprintf("%s is not an absolute path", source), nil, StatusError
		}
		if !filepath.IsAbs(destination) {
			return nil, fmt.Sprintf("%s is not an absolute path", destination), nil, StatusError
		}

		apply = func(out Output) bool {
			out.Runningf("Moving %s to %s", source, destination)
			if err := os.Rename(source, destination); err != nil {
				out.Errorf("Could not move %s to %s: %s", source, destination, err)
				return false
			}
			out.Successf("Moved %s to %s", source, destination)
			return true
		}
		return nil, fmt.Sprintf("Need to move %s to %s", source, destination), apply, StatusInfo
	},
	nil,
}
