package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

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
	extMeta := struct {
		Version int
	}{}
	jsondec.Decode(&extMeta)
	out.Info(fmt.Sprintf("Extension %s version is %d", ext, extMeta.Version))
	storeVersion(strconv.Itoa(extMeta.Version))
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
