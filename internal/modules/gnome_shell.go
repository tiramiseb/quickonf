package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("gnome-shell-extension", GnomeShellExtension)
	Register("local-gnome-shell-extension-version", LocalGnomeShellExtensionVersion)
	Register("gnome-shell-restart", GnomeShellRestart)
}

// GnomeShellExtension enables GNOME Shell extensions
func GnomeShellExtension(in interface{}, out output.Output) error {
	out.InstructionTitle("Enable GNOME Shell extension")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, extension := range data {
		if Dryrun {
			out.Info("Would enable " + extension)
			continue
		}
		if _, err := helper.Exec(nil, "gnome-shell-extension-tool", "--enable-extension", extension); err != nil {
			return err
		}
		out.Success("Enabled " + extension)
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
		return errors.New("Missing extension name")
	}
	f, err := os.Open(helper.Path(fmt.Sprintf(".local/share/gnome-shell/extensions/%s/metadata.json", ext)))

	storeVersion := func(val string) {
		if storeAs, ok := data["store"]; ok {
			helper.Store(storeAs, val)
		}
	}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			out.Info("Extension " + ext + " is not installed")
			storeVersion("")
			return nil
		}
		return err
	}
	defer f.Close()
	jsondec := json.NewDecoder(f)
	extMeta := map[string]interface{}{}
	jsondec.Decode(&extMeta)
	version, ok := extMeta["version"]
	if !ok {
		out.Info("Extension " + ext + " does not declare a version")
		storeVersion("")
		return nil
	}
	switch val := version.(type) {
	case string:
		out.Info(fmt.Sprintf("Extension %s version is %s", ext, val))
		storeVersion(val)
	case int:
		out.Info(fmt.Sprintf("Extension %s version is %d.0.0", ext, val))
		storeVersion(fmt.Sprintf("%d.0.0", val))
	case float64:
		out.Info(fmt.Sprintf("Extension %s version is %.f.0.0", ext, val))
		storeVersion(fmt.Sprintf("%.f.0.0", val))
	default:
		out.Alert(fmt.Sprintf("Extension %s version is not understood", ext))
		storeVersion("")
	}
	return nil
}

// GnomeShellRestart restarts GNOME Shell
func GnomeShellRestart(in interface{}, out output.Output) error {
	out.InstructionTitle("Restart GNOME Shell")
	if Dryrun {
		out.Info("Would restart GNOME Shell")
		return nil
	}
	if _, err := helper.Exec(nil, "busctl", "--user", "call", "org.gnome.Shell", "/org/gnome/Shell", "org.gnome.Shell", "Eval", "s", `Meta.restart("Restarting Gnome...")`); err != nil {
		return err
	}
	return nil
}
