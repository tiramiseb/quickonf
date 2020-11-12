package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("sudo-password", SudoPassword)
}

func SudoPassword(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Set password for sudo")
	data, err := input.String(in, store)
	if err != nil {
		return err
	}
	helper.SudoPassword = data
	out.Success("Password stored in memory (not displayed here, of course)")
	return nil
}
