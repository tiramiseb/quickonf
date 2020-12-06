package modules

import (
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("flatpak", Flatpak)
	Register("flatpak-remote", FlatpakRemote)
}

// Flatpak installs flatpack packages
func Flatpak(in interface{}, out output.Output) error {
	out.InstructionTitle("Installing flatpak package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		if Dryrun {
			out.Info("Would install " + pkg)
			continue
		}
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

// FlatpakRemote installs a flatpack package from a remote location
func FlatpakRemote(in interface{}, out output.Output) error {
	out.InstructionTitle("Adding flatpak repository")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for name, url := range data {
		if Dryrun {
			out.Info("Would add " + name + " (" + url + ")")
			continue
		}
		out.Info("Adding " + name + " (" + url + ")")
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
