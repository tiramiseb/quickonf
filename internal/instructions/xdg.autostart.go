package instructions

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/output"
)

const (
	xdgApplicationsDir = "/usr/share/applications"
	xdgAutostartDir    = "/etc/xdg/autostart/"
)

func init() {
	register(xdgAutostart)
}

var xdgAutostart = Instruction{
	"xdg.autostart",
	"Mark an application as auto-starting",
	"Do not change auto-start status",
	[]string{
		"Name of the application or path to the .desktop file",
	},
	nil,
	"Autostart GIMP\n  apt gimp\n  xdg.autostart gimp",
	func(args []string, out *output.Instruction, dry bool) ([]string, bool) {
		app := args[0]

		// Check the application .desktop file
		if !strings.HasSuffix(app, ".desktop") {
			app += ".desktop"
		}
		fullApp := app
		if !path.IsAbs(fullApp) {
			fullApp = filepath.Join(xdgApplicationsDir, fullApp)
			if _, err := os.Stat(fullApp); err != nil {
				out.Errorf("could not find %s", app)
				return nil, false
			}
		}

		// Check if the symlink already exists
		autostartPath := filepath.Join(xdgAutostartDir, filepath.Base(app))
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
		out.Successf("created %s", autostartPath)
		return nil, true
	},
}
