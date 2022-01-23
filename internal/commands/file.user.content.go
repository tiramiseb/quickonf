package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
)

func init() {
	register(fileUserContent)
}

var fileUserContent = Command{
	"file.user.content",
	"Create a file owned by a user (if path is relative, it is relative to the user's home directory)",
	[]string{
		"Username",
		"Path of the file",
		"Content of the file",
	},
	nil,
	"Say hello in /home/alice/hello.txt\n  file.user.content alice hello.txt \"Hello World!\"",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		path := args[1]
		content := args[2]

		usr, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join(usr.User.HomeDir, path)
		}

		finfo, err := os.Lstat(path)
		var (
			existingContent string
			ownershipOk     bool
		)
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

			if stat, ok := finfo.Sys().(*syscall.Stat_t); ok {
				ownershipOk = usr.Uid == int(stat.Uid)
			}
		}

		var willMessage string
		var needMessage string
		switch {
		case content == existingContent && ownershipOk:
			return nil, fmt.Sprintf("%s already has the requested content", path), nil, StatusSuccess
		case content == existingContent && !ownershipOk:
			willMessage = fmt.Sprintf("Will change ownership of %s to %s", path, username)
			needMessage = fmt.Sprintf("Need to change ownership of %s to %s", path, username)
		default:
			willMessage = fmt.Sprintf("Will write requested content to %s", path)
			needMessage = fmt.Sprintf("Need to write requested content to %s", path)
		}

		apply = &Apply{
			"file.user.content",
			willMessage,
			func(out Output) bool {
				if existingContent != content {
					out.Infof("Writing content to %s", path)
					if err := os.WriteFile(path, []byte(content), 0644); err != nil {
						out.Errorf("Could not write requested content to %s: %s", path, err)
						return false
					}
					if !ownershipOk {
						out.Infof("Content written to %s", path)
					}
				}
				if !ownershipOk {
					out.Infof("Changing ownership of %s", path)
					if err := os.Chown(path, usr.Uid, usr.Group.Gid); err != nil {
						out.Errorf("Could not change ownership of %s: %s", path, err)
						return false
					}
					if existingContent == content {
						out.Successf("Changed ownership of %s", path)
						return true
					}
				}
				out.Successf("Content written to %s", path)
				return true
			},
		}

		return nil, needMessage, apply, StatusInfo
	},
	datastores.Users.Reset,
}
