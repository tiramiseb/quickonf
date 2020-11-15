package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("dconf", Dconf)
}

// Dconf sets a parameter in the dconf database
func Dconf(in interface{}, out output.Output) error {
	out.InstructionTitle("Dconf database")
	data, err := helper.MapStringString(in)
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