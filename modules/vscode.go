package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("vscode-extension", VscodeExtension)
	Register("vscodium-extension", VscodiumExtension)
}

func VscodeExtension(in interface{}, out output.Output, store map[string]interface{}) error {
	return vscodeExtension(in, out, "code", store)
}

func VscodiumExtension(in interface{}, out output.Output, store map[string]interface{}) error {
	return vscodeExtension(in, out, "codium", store)
}

func vscodeExtension(in interface{}, out output.Output, cmd string, store map[string]interface{}) error {
	out.ModuleName("Install Vs" + cmd + " extension")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, extension := range data {
		out.Info("Installing " + extension)
		out.ShowLoader()
		_, err := helper.Exec(cmd, "--install-extension", extension)
		out.HideLoader()
		if err != nil {
			return nil
		}
	}
	return nil
}
