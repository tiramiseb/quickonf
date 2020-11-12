package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("update-alternatives", UpdateAlternatives)
}

func UpdateAlternatives(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Update Alternative")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for alt, path := range data {
		err := helper.ExecSudo("update-alternatives", "--set", alt, path)
		if err != nil {
			return err
		}
		out.Success("Changed alternative for " + alt + " to " + path)
	}
	return nil
}
