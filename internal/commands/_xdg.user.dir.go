package commands

import (
	"bytes"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
)

var xdgDirs = []string{
	"DESKTOP",
	"DOWNLOAD",
	"TEMPLATES",
	"PUBLICSHARE",
	"DOCUMENTS",
	"MUSIC",
	"PICTURES",
	"VIDEOS",
}

func isXdgDir(dir string) bool {
	for _, d := range xdgDirs {
		if d == dir {
			return true
		}
	}
	return false
}

var xdgDirsList = strings.Join(xdgDirs, ", ")

func init() {
	register(xdgUserDir)
}

var xdgUserDir = Command{
	"xdg.user.dir",
	"Set a XDG user dir",
	"Do not change XDG user dir",
	[]string{
		"Username",
		"Directory name (one of " + xdgDirsList + ", case-insensitive)",
		"Directory path",
	},
	nil,
	"Change downloads directory for john\n  file.user.directory Downs\n  xdg.user.dir john DOWNLOAD Downs",
	func(args []string, out output, dry bool) ([]string, bool) {
		username := args[0]
		name := args[1]
		path := args[2]

		usr, err := user.Lookup(username)
		if err != nil {
			out.Errorf("could not identify user %s: %s", username, err)
			return nil, false
		}

		name = strings.ToUpper(name)
		if !isXdgDir(name) {
			out.Errorf("XDG user dir %s is invalid", name)
			return nil, false
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join(usr.HomeDir, path)
		}

		var buf bytes.Buffer
		wait, err := helper.ExecAs(usr, nil, &buf, "xdg-user-dir", name)
		if err != nil {
			out.Errorf("could not check XDG user dir %s: %s (%s)", name, err, buf.String())
			return nil, false
		}
		if err := wait(); err != nil {
			out.Errorf("could not check XDG user dir %s: %s", name, buf.String())
			return nil, false
		}

		current := strings.TrimSpace(buf.String())
		if current == path {
			out.Successf("XDG user dir %s for %s is already %s", name, username, path)
			return nil, true
		}

		if dry {
			out.Infof("would set XDG user dir %s for %s to %s", name, username, path)
			return nil, true
		}

		buf.Reset()
		wait, err = helper.ExecAs(usr, nil, &buf, "xdg-user-dirs-update", "--set", name, path)
		if err != nil {
			out.Errorf("could not change XDG user dir %s: %s", name, err)
			return nil, false
		}
		if err := wait(); err != nil {
			out.Errorf("could not change XDG user dir %s: %s", name, buf.String())
			return nil, false
		}
		out.Successf("XDG user dir %s for %s set to %s", name, username, path)
		return nil, true
	},
}
