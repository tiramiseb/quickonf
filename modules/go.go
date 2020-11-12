package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("go-env", GoEnv)
	Register("go-package", GoPackage)
}

func GoEnv(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Add go environment variable")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for param, value := range data {
		if _, err := helper.Exec("go", "env", "-w", param+"="+value); err != nil {
			return err
		}
		out.Success("Set " + param + " to " + value)
	}
	return nil
}

func GoPackage(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Install go package")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, pkg := range data {
		out.Info("Installing " + pkg)
		out.ShowLoader()
		_, err := helper.Exec("go", "get", pkg)
		out.HideLoader()
		if err != nil {
			return err
		}
	}
	return nil
}
