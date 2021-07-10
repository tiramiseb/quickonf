package modules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("install-gnome-shell-extension", InstallGnomeShellExtension)
	Register("enable-gnome-shell-extension", EnableGnomeShellExtension)
	Register("local-gnome-shell-extension-version", LocalGnomeShellExtensionVersion)
	Register("gnome-shell-restart", GnomeShellRestart)
}

// InstallGnomeShellExtension installs GNOME Shell extensions from extensions.gnome.org
func InstallGnomeShellExtension(in interface{}, out output.Output) error {
	out.InstructionTitle("Install GNOME Shell extension")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	cmdout, err := helper.Exec(nil, "", "gnome-shell", "--version")
	if err != nil {
		return err
	}
	gnomeVerLine := bytes.Fields(cmdout)
	if len(gnomeVerLine) != 3 {
		return fmt.Errorf(`GNOME version invalid, should be "GNOME Shell X.Y.Z": %s`, cmdout)
	}
	gnomeVersion := string(gnomeVerLine[2])
	for _, ext := range data {
		dest := helper.Path(".local/share/gnome-shell/extensions/" + ext)
		out.ShowLoader()
		extInfo := struct {
			Version     int
			DownloadURL string `json:"download_url"`
		}{}
		err = helper.DownloadJSON("https://extensions.gnome.org/extension-info/?uuid="+ext+"&shell_version="+gnomeVersion, &extInfo)
		out.HideLoader()
		if err != nil {
			if err.Error() == "404 not found" {
				return fmt.Errorf("extension %s does not exist", ext)
			}
			return err
		}

		_, current, err := localGnomeShellExtensionVersion(ext)
		if err != nil {
			return err
		}
		if extInfo.Version <= current {
			out.Infof("%s already installed in version %d", ext, current)
			continue
		}
		if Dryrun {
			out.Infof("Would install %s in version %d", ext, extInfo.Version)
			continue
		}

		out.Infof("Installing %s", ext)
		tmpfile, err := ioutil.TempFile("", "quickonf-gnome-extension-"+ext+"-*.zip")
		if err != nil {
			return err
		}
		fname := tmpfile.Name()
		if err := tmpfile.Close(); err != nil {
			return err
		}
		if _, err := helper.DownloadFile("https://extensions.gnome.org"+extInfo.DownloadURL, fname, out); err != nil {
			return err
		}
		if _, err := helper.ExtractZip(fname, dest, out); err != nil {
			return err
		}
	}
	return nil
}

// EnableGnomeShellExtension enables GNOME Shell extensions
func EnableGnomeShellExtension(in interface{}, out output.Output) error {
	out.InstructionTitle("Enable GNOME Shell extension")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, extension := range data {
		if Dryrun {
			out.Infof("Would enable %s", extension)
			continue
		}
		if _, err := helper.Exec(nil, "", "gnome-shell-extension-tool", "--enable-extension", extension); err != nil {
			return err
		}
		out.Successf("Enabled %s", extension)
	}
	return nil
}

// LocalGnomeShellExtensionVersion checks a locally installed GNOME Shell extension version
func LocalGnomeShellExtensionVersion(in interface{}, out output.Output) error {
	out.InstructionTitle("Checking GNOME Shell extension version")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	ext, ok := data["extension"]
	if !ok {
		return errors.New("missing extension name")
	}

	installed, version, err := localGnomeShellExtensionVersion(ext)
	if err != nil {
		return err
	}

	if installed {
		out.Infof("Extension %s current version is %d", ext, version)
	} else {
		out.Infof("Extension %s is not installed", ext)
	}
	helper.Store(data["store"], strconv.Itoa(version))
	return nil
}

func localGnomeShellExtensionVersion(ext string) (bool, int, error) {
	f, err := os.Open(helper.Path(fmt.Sprintf(".local/share/gnome-shell/extensions/%s/metadata.json", ext)))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return false, 0, err
	}
	defer f.Close()
	jsondec := json.NewDecoder(f)
	extMeta := struct {
		Version int
	}{}
	if err := jsondec.Decode(&extMeta); err != nil {
		return false, 0, err
	}
	return true, extMeta.Version, nil

}

// GnomeShellRestart restarts GNOME Shell
func GnomeShellRestart(in interface{}, out output.Output) error {
	out.InstructionTitle("Restart GNOME Shell")
	if Dryrun {
		out.Info("Would restart GNOME Shell")
		return nil
	}
	if _, err := helper.Exec(nil, "", "busctl", "--user", "call", "org.gnome.Shell", "/org/gnome/Shell", "org.gnome.Shell", "Eval", "s", `Meta.restart("Restarting Gnome...")`); err != nil {
		return err
	}
	return nil
}
