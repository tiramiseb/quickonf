package modules

import (
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("pip", PIP)
}

// PIP installs Python packages using pip
func PIP(in interface{}, out output.Output) error {
	out.InstructionTitle("Installing PIP package")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		if Dryrun {
			out.Info("Would install " + pkg)
			continue
		}
		if err = helper.ExecSudo("pip3", "install", "--upgrade", pkg); err != nil {
			if strings.Contains(err.Error(), "command not found") {
				out.Info("Installing PIP first")
				out.ShowLoader()
				err = helper.ExecSudo("apt-get", "--yes", "--quiet", "install", "--no-install-recommends", "python3-pip")
				out.HideLoader()
				if err != nil {
					return err
				}
				if err := helper.ExecSudo("pip3", "install", "--upgrade", pkg); err != nil {
					return err
				}
			}
		}
	}
	return err
}
