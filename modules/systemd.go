package modules

import (
	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("systemd-enable", SystemdEnable)
}

func SystemdEnable(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Enabling service")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, service := range data {
		err = helper.ExecSudo("systemctl", "enable", service)
		if err != nil {
			return err
		}
		err = helper.ExecSudo("systemctl", "start", service)
		if err != nil {
			return err
		}
        out.Success("Enabled "+service)
	}
	return nil
}
