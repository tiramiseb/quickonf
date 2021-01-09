package modules

import (
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
			out.Infof("Would install %s", pkg)
			continue
		}
		if _, err := helper.ExecSudo(nil, "", "pip3", "install", "--upgrade", pkg); err != nil {
			return err
		}
	}
	return nil
}
