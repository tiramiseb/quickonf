package modules

import (
	"syscall"

	"github.com/shirou/gopsutil/process"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("gnome-shell-extension", GnomeShellExtension)
	Register("gnome-shell-restart", GnomeShellRestart)
}

func GnomeShellExtension(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Enable GNOME Shell extension")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, extension := range data {
		if _, err := helper.Exec("gnome-shell-extension-tool", "--enable-extension", extension); err != nil {
			return err
		}
		out.Success("Enabled " + extension)
	}
	return nil
}

func GnomeShellRestart(in interface{}, out output.Output, store map[string]interface{}) error {
    // TODO busctl --user call org.gnome.Shell /org/gnome/Shell org.gnome.Shell Eval s 'Meta.restart("Restarting Gnome...")';
	out.ModuleName("Restart GNOME Shell")
	procs, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range procs {
		n, err := p.Name()
		if err != nil {
			return err
		}
		if n == "gnome-shell" {
			p.SendSignal(syscall.SIGQUIT)
		}
	}
	return nil
}
