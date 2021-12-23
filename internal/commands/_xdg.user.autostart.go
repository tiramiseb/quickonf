package commands

import (
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	xdgUserApplicationsDir = ".local/share/applications"
	xdgUserAutostartDir    = ".config/autostart/"
)

func init() {
	register(xdgUserAutostart)
}

var xdgUserAutostart = Command{
	"xdg.user.autostart",
	"Mark an application as auto-starting, for a specific user",
	"Do not change auto-start status",
	[]string{
		"Username",
		"Name of the application or path to the .desktop file",
	},
	nil,
	"Autostart GIMP\n  apt gimp\n  xdg.user.autostart john gimp",
	func(args []string, out output, dry bool) ([]string, bool) {
		username := args[0]
		app := args[1]

		usr, err := user.Lookup(username)
		if err != nil {
			out.Errorf("could not identify user %s: %v", username, err)
			return nil, false
		}

		// Check the application .desktop file
		if !strings.HasSuffix(app, ".desktop") {
			app += ".desktop"
		}
		fullApp := app
		if !path.IsAbs(fullApp) {
			fullApp = filepath.Join(xdgApplicationsDir, fullApp)
			if _, err := os.Stat(fullApp); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					out.Errorf("could not find %s", app)
					return nil, false
				}
				fullApp = filepath.Join(usr.HomeDir, xdgUserApplicationsDir, fullApp)
				if _, err := os.Stat(fullApp); err != nil {
					out.Errorf("could not find %s", app)
					return nil, false
				}
			}
		}

		// Check if the symlink already exists
		autostartPath := filepath.Join(usr.HomeDir, xdgUserAutostartDir, filepath.Base(app))
		if stat, err := os.Lstat(autostartPath); err == nil {
			if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
				target, err2 := os.Readlink(autostartPath)
				if err2 != nil {
					out.Errorf("could not check already-existing %s: %v", autostartPath, err2)
					return nil, false
				}
				if target == fullApp {
					out.Successf("%s already auto-started", app)
					return nil, true
				}
			}
			if dry {
				out.Infof("would remove %s", autostartPath)
			} else {
				if err2 := os.Remove(autostartPath); err2 != nil {
					out.Errorf("could not remove %s: %v", autostartPath, err2)
					return nil, false
				}
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			out.Errorf("could not check already-existing %s: %v", autostartPath, err)
			return nil, false
		}

		// Create the symlink
		if dry {
			out.Infof("would create the %s symlink", autostartPath)
			return nil, true
		}
		if err := os.Symlink(fullApp, autostartPath); err != nil {
			out.Errorf("could not create %s: %v", autostartPath, err)
			return nil, false
		}
		uid, err := strconv.Atoi(usr.Uid)
		if err != nil {
			out.Errorf("could not get user ID: %v", err)
			return nil, false
		}
		gid, err := strconv.Atoi(usr.Gid)
		if err != nil {
			out.Errorf("could not get group ID: %v", err)
			return nil, false
		}
		if err := os.Chown(autostartPath, uid, gid); err != nil {
			out.Errorf("could not change ownership of %s: %v", autostartPath, err)
			return nil, false
		}
		out.Successf("created %s", autostartPath)
		return nil, true
	},
}
