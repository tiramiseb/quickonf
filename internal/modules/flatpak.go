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
			out.Infof("Would install %s", pkg)
			continue
		}
		out.Infof("Installing %s", pkg)
		out.ShowLoader()
		_, err := helper.ExecSudo(nil, "", "flatpak", "install", "--noninteractive", "--assumeyes", pkg)
		out.HideLoader()
		if err != nil {
			if strings.Contains(err.Error(), "already installed") {
				out.Info("... already installed")
				continue
			}
			return err
		}
	}
	return nil
}

// FlatpakRemote adds a flatpak remote repository
func FlatpakRemote(in interface{}, out output.Output) error {
	out.InstructionTitle("Adding flatpak repository")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for name, url := range data {
		if Dryrun {
			out.Infof("Would add %s (%s)", name, url)
			continue
		}
		out.Infof("Adding %s (%s)", name, url)
		out.ShowLoader()
		_, err := helper.ExecSudo(nil, "", "flatpak", "remote-add", "--if-not-exists", name, url)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	out.Info("You may need to reboot after flatpak configuration")
	return nil
}
