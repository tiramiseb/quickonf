package modules

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

var (
	xdgApplicationsDir = "/usr/share/applications"
	xdgAutostartDir    = ".config/autostart"
	xdgAllUserDirs     = map[string]bool{
		"DESKTOP":     true,
		"DOCUMENTS":   true,
		"DOWNLOAD":    true,
		"MUSIC":       true,
		"PICTURES":    true,
		"PUBLICSHARE": true,
		"TEMPLATES":   true,
		"VIDEOS":      true,
	}
)

func init() {
	Register("xdg-autostart", XdgAutostart)
	Register("xdg-mime-default", XdgMimeDefault)
	Register("xdg-user-dir", XdgUserDir)
}

// XdgAutostart enables autostart for given applications
func XdgAutostart(in interface{}, out output.Output) error {
	out.InstructionTitle("XDG autostart")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, app := range data {
		appBase := app
		if !strings.HasSuffix(appBase, ".desktop") {
			appBase = app + ".desktop"
		}
		fullApp := appBase
		if !path.IsAbs(fullApp) {
			if _, err := os.Stat(filepath.Join(xdgApplicationsDir, fullApp)); !errors.Is(err, os.ErrNotExist) {
				fullApp = filepath.Join(xdgApplicationsDir, fullApp)
			}
		}
		fullApp = helper.Path(fullApp)
		autostartPath := helper.Path(filepath.Join(xdgAutostartDir, filepath.Base(appBase)))

		status, err := helper.Symlink(autostartPath, fullApp)
		switch status {
		case helper.SymlinkError:
			return err
		case helper.SymlinkAleradyExists:
			out.Info("Autostart for " + app + " already configured")
		case helper.SymlinkCreated:
			if Dryrun {
				out.Info("Would autostart " + app)
			} else {
				out.Success("Autostart for " + app + " configured")
			}
		}
	}
	return nil
}

// XdgMimeDefault sets default applications for mimetypes
func XdgMimeDefault(in interface{}, out output.Output) error {
	out.InstructionTitle("XDG mime default")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for mimetype, app := range data {
		if Dryrun {
			out.Info("Would change default app for " + mimetype + " to " + app)
			continue
		}
		if _, err := helper.Exec("xdg-mime", "default", app+".desktop", mimetype); err != nil {
			return err
		}
		out.Success("Changed default app for " + mimetype + " to " + app)
	}
	return nil
}

// XdgUserDir sets XDG user dirs
func XdgUserDir(in interface{}, out output.Output) error {
	out.InstructionTitle("XDG user dir")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for name, path := range data {
		name = strings.ToUpper(name)
		if !xdgAllUserDirs[name] {
			return errors.New("User dir \"" + name + "\" does not exist")
		}
		path = helper.Path(path)
		if Dryrun {
			out.Info("Would change user dir " + name + " to " + path)
			continue
		}
		if _, err := helper.Exec("xdg-user-dirs-update", "--set", name, path); err != nil {
			return err
		}
		out.Success("Changed user dir " + name + " to " + path)
	}
	return nil
}
