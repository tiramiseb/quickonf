package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/tiramiseb/quickonf/commands/datastores"
)

func init() {
	register(userFileLine)
}

var userFileLine = &Command{
	"user.file.line",
	"Make sure a file, owned by a user, contains a line (if path is relative, it is relative to the user's home directory)",
	[]string{
		"Username",
		"Path of the file",
		"Line that must be in the file",
		"Regexp for line(s) to replace",
	},
	nil,
	"Have hello\n  user.file.line alice hello.txt \"hello Alice!\" \"^hello\"",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		username := args[0]
		path := args[1]
		line := args[2]
		matcher := args[3]

		usr, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(usr.User.HomeDir, path)
		}

		re, err := regexp.Compile(matcher)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		var existingLines []string
		finfo, err := os.Lstat(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err.Error(), nil, StatusError, "", ""
			}
		} else {
			if finfo.IsDir() {
				return nil, fmt.Sprintf("%s is a directory", path), nil, StatusError, "", ""
			}
			f, err := os.Open(path)
			if err != nil {
				return nil, err.Error(), nil, StatusError, "", ""
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				thisLine := scanner.Bytes()
				if re.Match(thisLine) {
					existingLines = append(existingLines, string(thisLine))
				}
			}
			if err := scanner.Err(); err != nil {
				return nil, err.Error(), nil, StatusError, "", ""
			}
		}

		if len(existingLines) == 1 && existingLines[0] == line {
			return nil, fmt.Sprintf("%s already has the requested line", path), nil, StatusSuccess, line, line
		}
		apply = func(out Output) bool {
			out.Runningf("Writing line to %s", path)
			finfo, err := os.Lstat(path)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					out.Errorf("Could not check %s: %s", path, err)
					return false
				}
				if err := os.WriteFile(path, []byte(line), 0o644); err != nil {
					out.Errorf("Could not write requested content to %s: %s", path, err)
					return false
				}
				out.Successf("Content written to %s", path)
				return true
			}
			src, err := os.Open(path)
			if err != nil {
				out.Errorf("Could not open %s: %s", path, err)
				return false
			}
			var buf bytes.Buffer
			wroteIt := false
			scanner := bufio.NewScanner(src)
			for scanner.Scan() {
				thisLine := scanner.Bytes()
				if re.Match(thisLine) {
					if !wroteIt {
						buf.WriteString(line)
						buf.WriteByte('\n')
						wroteIt = true
					}
				} else {
					buf.Write(thisLine)
					buf.WriteByte('\n')
				}
			}
			src.Close()
			if err := scanner.Err(); err != nil {
				out.Errorf("Could not read %s: %s", path, err)
				return false
			}
			if err := os.WriteFile(path, buf.Bytes(), finfo.Mode()); err != nil {
				out.Errorf("Could not write line to %s: %s", path, err)
				return false
			}
			if stat, ok := finfo.Sys().(*syscall.Stat_t); ok {
				if usr.Uid != int(stat.Uid) {
					out.Runningf("Changing ownership of %s", path)
					if err := os.Chown(path, usr.Uid, usr.Group.Gid); err != nil {
						out.Errorf("Could not change ownership of %s: %s", path, err)
						return false
					}
				}
			}
			out.Successf("Content written to %s", path)
			return true
		}

		return nil, fmt.Sprintf("Need to write requested content to %s", path), apply, StatusInfo, strings.Join(existingLines, "\n[...]\n"), line
	},
	nil,
}
