package modules

import (
	"errors"
	"fmt"
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
			out.Infof("Autostart for %s already configured", app)
		case helper.SymlinkCreated:
			if Dryrun {
				out.Infof("Would autostart %s", app)
			} else {
				out.Successf("Autostart for %s configured", app)
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
			out.Infof("Would change default app for %s to %s", mimetype, app)
			continue
		}
		if _, err := helper.Exec(nil, "xdg-mime", "default", app+".desktop", mimetype); err != nil {
			return err
		}
		out.Successf("Changed default app for %s to %s", mimetype, app)
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
			return fmt.Errorf(`User dir "%s" does not exist`, name)
		}
		path = helper.Path(path)
		if Dryrun {
			out.Infof("Would change user dir %s to %s", name, path)
			continue
		}
		if _, err := helper.Exec(nil, "xdg-user-dirs-update", "--set", name, path); err != nil {
			return err
		}
		out.Successf("Changed user dir %s to %s", name, path)
	}
	return nil
}
