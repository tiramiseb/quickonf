package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("systemd-enable", SystemdEnable)
	Register("systemd-disable", SystemdDisable)
}

// SystemdEnable enables a systemd service
func SystemdEnable(in interface{}, out output.Output) error {
	out.InstructionTitle("Enabling service")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, service := range data {
		if Dryrun {
			out.Info("Would enable and start " + service)
			continue
		}
		err = helper.ExecSudo("systemctl", "enable", service)
		if err != nil {
			return err
		}
		err = helper.ExecSudo("systemctl", "start", service)
		if err != nil {
			return err
		}
		out.Success("Enabled " + service)
	}
	return nil
}

// SystemdDisable disables a systemd service
func SystemdDisable(in interface{}, out output.Output) error {
	out.InstructionTitle("Disabling service")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, service := range data {
		if Dryrun {
			out.Info("Would stop and disable " + service)
			continue
		}
		err = helper.ExecSudo("systemctl", "stop", service)
		if err != nil {
			return err
		}
		err = helper.ExecSudo("systemctl", "disable", service)
		if err != nil {
			return err
		}
		out.Success("Disabled " + service)
	}
	return nil
}
