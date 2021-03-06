package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("go-env", GoEnv)
	Register("go-package", GoPackage)
}

// GoEnv sets go environment parameters
func GoEnv(in interface{}, out output.Output) error {
	out.InstructionTitle("Add go environment variable")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for param, value := range data {
		if Dryrun {
			out.Infof("Would set %s to %s", param, value)
			continue
		}
		if _, err := helper.Exec(nil, "", "go", "env", "-w", param+"="+value); err != nil {
			return err
		}
		out.Successf("Set %s to %s", param, value)
	}
	return nil
}

// GoPackage install go packages
func GoPackage(in interface{}, out output.Output) error {
	out.InstructionTitle("Install go package")
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
		_, err := helper.Exec(nil, "", "go", "get", pkg)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	return nil
}
