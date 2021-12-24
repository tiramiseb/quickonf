package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
	"github.com/tiramiseb/quickonf/internal/commands/shared"
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

func resetXdgDir() {
	xdgDirsStorage = map[string]map[string]string{}
}

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
			dirs[matches[1]] = matches[2]
		}
		xdgDirsStorage[usr.Username] = dirs
	}
	log.Print("USER DIRS")
	log.Print(dirs)
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
	"Change downloads directory for john\n  file.user.directory Downs\n  xdg.user.dir john DOWNLOAD Downs",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		username := args[0]
		name := args[1]
		path := args[2]

		usr, err := shared.Users.Get(username)
		if err != nil {
			return nil, err.Error(), nil, StatusError
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join("$HOME", path)
		}

		dir, err := xdgDir(usr, name)
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
				var buf bytes.Buffer
				wait, err := helper.ExecAs(usr, nil, &buf, "xdg-user-dirs-update", "--set", name, path)
				if err != nil {
					out.Errorf("could not change XDG user dir %s: %s", name, err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("could not change XDG user dir %s: %s", name, buf.String())
					return false
				}
				out.Successf("XDG user dir %s for %s set to %s", name, username, path)
				return true
			},
		}

		return nil, fmt.Sprintf("need to set %s directory to %s", name, path), apply, StatusInfo
	},
	resetXdgDir,
}
