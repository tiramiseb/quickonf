package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("adduser", Adduser)
}

func Adduser(in interface{}, out output.Output, store map[string]interface{}) error {
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	for user, group := range data {
		helper.ExecSudo("adduser", user, group)
		if err != nil {
			return err
		}
		out.Success("User " + user + " added to group " + group)
	}
	return nil
}
