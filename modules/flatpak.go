package modules

import (
    "strings"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("flatpak", Flatpak)
	Register("flatpak-remote", FlatpakRemote)
}

func Flatpak(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Installing flatpak package")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		out.Info("Installing " + pkg)
		out.ShowLoader()
        err := helper.ExecSudo("flatpak", "install", "--noninteractive", "--assumeyes", pkg)
		out.HideLoader()
		if err != nil && strings.Index(err.Error(), "already installed") == -1 {
			return err
		}
	}
	return nil
}

func FlatpakRemote(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Installing flatpak package")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for name, url := range data {
        out.Info("Adding "+name+" ("+url+")")
		out.ShowLoader()
        err := helper.ExecSudo("flatpak", "remote-add", "--if-not-exists", name, url)
		out.HideLoader()
        if err != nil {
            return err
        }
    }
    out.Info("You may need to reboot after flatpak configuration")
    return nil
}
