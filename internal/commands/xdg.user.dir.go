package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

const xdgUserDirsSubpath = ".config/user-dirs.dirs"

var xdgDirsValidNames = []string{
	"DESKTOP",
	"DOWNLOAD",
	"TEMPLATES",
	"PUBLICSHARE",
	"DOCUMENTS",
	"MUSIC",
	"PICTURES",
	"VIDEOS",
}

var (
	xdgDirsStorage = map[string]map[string]string{}
	xdgDirRe       = regexp.MustCompile(`^XDG_(.*)_DIR="(.*)"`)
)

func xdgDir(usr *user.User, dirname string) (string, error) {
	dirs, ok := xdgDirsStorage[usr.Username]
	if !ok {
		f, err := os.Open(filepath.Join(usr.HomeDir, xdgUserDirsSubpath))
		if err != nil {
			return "", err
		}
		dirs = make(map[string]string, len(xdgDirsValidNames))
		for _, k := range xdgDirsValidNames {
			dirs[k] = ""
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			matches := xdgDirRe.FindStringSubmatch(line)
			if len(matches) != 3 {
				continue
			}
			dirs[matches[1]] = strings.Replace(matches[2], "$HOME", usr.HomeDir, 1)
		}
		xdgDirsStorage[usr.Username] = dirs
	}
	dir, ok := dirs[strings.ToUpper(dirname)]
	if !ok {
		return "", fmt.Errorf(`"%s" is not a valid XDG user dir`, dirname)
	}
	return dir, nil
}

func init() {
	register(xdgUserDir)
}

var xdgUserDir = Command{
	"xdg.user.dir",
	"Set a XDG user dir",
	[]string{
		"Username",
		"Directory name (one of " + strings.Join(xdgDirsValidNames, ", ") + ", case-insensitive)",
		"Directory path",
	},
	nil,
	"Change downloads directory for john\n  file.user.directory john Downs\n  xdg.user.dir john DOWNLOAD Downs",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		name := strings.ToUpper(args[1])
		path := args[2]

		usr, err := datastores.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if path == "" {
			path = usr.User.HomeDir + "/"
		} else if !filepath.IsAbs(path) {
			path = filepath.Join(usr.User.HomeDir, path)
		}

		dir, err := xdgDir(usr.User, name)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if dir == path {
			return nil, fmt.Sprintf("%s is already %s", name, dir), nil, StatusSuccess
		}

		apply = &Apply{
			"xdg.user.dir",
			fmt.Sprintf("Will set %s directory to %s", name, path),
			func(out Output) bool {
				out.Infof("Setting %s directory to %s", name, path)
				if err := helper.ExecAs(usr.User, nil, nil, "xdg-user-dirs-update", "--set", name, path); err != nil {
					out.Errorf("Could not change XDG user dir %s: %s", name, helper.ExecErr(err))
					return false
				}
				out.Successf("XDG user dir %s for %s set to %s", name, username, path)
				return true
			},
		}

		return nil, fmt.Sprintf("Need to set %s directory to %s", name, path), apply, StatusInfo
	},
	func() {
		xdgDirsStorage = map[string]map[string]string{}
		datastores.Users.Reset()
	},
}
