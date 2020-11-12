package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("git-config", GitConfig)
}

func GitConfig(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Git configuration")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for param, value := range data {
		if _, err := helper.Exec("git", "config", "--global", param, value); err != nil {
			return err
		}
		out.Success("Set " + param + " to " + value)
	}
	return nil
}
