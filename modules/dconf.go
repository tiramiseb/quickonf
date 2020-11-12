package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("dconf", Dconf)
}

func Dconf(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Dconf database")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for k, v := range data {
		_, err := helper.Exec("dconf", "write", k, v)
		if err != nil {
			return err
		}
		out.Success("Set " + k + " to " + v)
	}
	return nil
}
