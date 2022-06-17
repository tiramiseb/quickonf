package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func init() {
	register(filePerm)
}

var filePermRe = regexp.MustCompile("^[0-7][0-7][0-7]$")

var filePerm = Command{
	"file.perm",
	"Change a file permissions",
	[]string{
		"Absolute path of the file",
		"Permissions (000-777)",
	},
	nil,
	"Restrict private key\n  file.chown /home/alice/.ssh/id_rsa 600",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		path := args[0]
		permsAsString := args[1]
		if !filepath.IsAbs(path) {
			return nil, fmt.Sprintf("%s is not an absolute path", path), nil, StatusError, "", ""
		}
		if !filePermRe.MatchString(permsAsString) {
			return nil, fmt.Sprintf("%s is not a correct permission (000-777)", permsAsString), nil, StatusError, "", ""
		}
		permsAsInt, _ := strconv.ParseInt(permsAsString, 8, 64) // Not checking the error because it cannot fail, because mode has been checked by the regex
		perms := fs.FileMode(permsAsInt)

		finfo, err := os.Lstat(path)
		var currentPerms fs.FileMode
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err.Error(), nil, StatusError, "", ""
			}
		} else {
			currentPerms = finfo.Mode().Perm()
		}
		before = currentPerms.String()
		after = perms.String()
		if perms == currentPerms {
			return nil, fmt.Sprintf("%s already has the requested permissions", path), nil, StatusSuccess, before, after
		}

		apply = func(out Output) bool {
			out.Runningf("Changing permissions on %s", path)
			if err := os.Chmod(path, perms); err != nil {
				out.Errorf("Could not change permissions on %s: %s", path, err)
				return false
			}
			out.Successf("Changed permissions on %s", path)
			return true
		}

		return nil, fmt.Sprintf("Need to change permissions on %s", path), apply, StatusInfo, before, after
	},
	nil,
}
