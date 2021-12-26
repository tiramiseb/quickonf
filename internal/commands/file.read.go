package commands

import (
	"fmt"
	"os"
)

func init() {
	register(fileRead)
}

var fileRead = Command{
	"file.read",
	"Read the content of a file",
	[]string{
		"Path of the file",
	},
	[]string{
		"Content of the file",
	},
	"APT sources\n  aptsrc = file.read <confdir>/sources.list\n  file.content /etc/apt/sources.list <aptsrc>",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		path := args[0]
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}
		return []string{string(content)}, fmt.Sprintf("Read content of file %s", path), nil, StatusInfo
	},
	nil,
}
