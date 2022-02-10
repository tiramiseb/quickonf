package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func init() {
	register(file)
}

var file = Command{
	"file.content",
	"Create a file owned by root",
	[]string{
		"Absolute path of the file",
		"Content of the file",
	},
	nil,
	"Say hello\n  file.content /home/hello.txt \"Hello World!\"",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]
		content := args[1]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError
		}
		finfo, err := os.Lstat(path)
		var existingContent string
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err.Error(), nil, StatusError
			}
		} else {
			if finfo.IsDir() {
				return nil, fmt.Sprintf("%s is a directory", path), nil, StatusError
			}
			bcontent, err := os.ReadFile(path)
			if err != nil {
				return nil, err.Error(), nil, StatusError
			}
			existingContent = string(bcontent)
		}
		if content == existingContent {
			return nil, fmt.Sprintf("%s already has the requested content", path), nil, StatusSuccess
		}

		apply = &Apply{
			"file.content",
			fmt.Sprintf("Will write requested content to %s", path),
			func(out Output) bool {
				out.Runningf("Writing content to %s", path)
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					out.Errorf("Could not write requested content to %s: %s", path, err)
					return false
				}
				out.Successf("Content written to %s", path)
				return true
			},
		}

		return nil, fmt.Sprintf("Need to write requested content to %s", path), apply, StatusInfo
	},
	nil,
}
