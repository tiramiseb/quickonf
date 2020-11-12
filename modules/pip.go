package modules

import (
	"strings"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("pip", PIP)
}

func PIP(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Installing PIP package")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		if err = helper.ExecSudo("pip3", "install", "--upgrade", pkg); err != nil {
			if strings.Contains(err.Error(), "command not found") {
				// Install pip if it is not already installed
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
