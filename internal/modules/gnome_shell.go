package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("gnome-shell-extension", GnomeShellExtension)
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
		if _, err := helper.Exec("gnome-shell-extension-tool", "--enable-extension", extension); err != nil {
			return err
		}
		out.Success("Enabled " + extension)
	}
	return nil
}

// GnomeShellRestart restarts GNOME Shell
func GnomeShellRestart(in interface{}, out output.Output) error {
	out.InstructionTitle("Restart GNOME Shell")
	if _, err := helper.Exec("busctl", "--user", "call", "org.gnome.Shell", "/org/gnome/Shell", "org.gnome.Shell", "Eval", "s", `Meta.restart("Restarting Gnome...")`); err != nil {
		return err
	}
	// procs, err := process.Processes()
	// if err != nil {
	// 	return err
	// }
	// for _, p := range procs {
	// 	n, err := p.Name()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if n == "gnome-shell" {
	// 		p.SendSignal(syscall.SIGQUIT)
	// 	}
	// }
	return nil
}
