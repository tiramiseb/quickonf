package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("systemd-enable", SystemdEnable)
	Register("systemd-disable", SystemdDisable)
	Register("systemd-user-enable", SystemdUserEnable)
	Register("systemd-user-disable", SystemdUserDisable)
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
			out.Infof("Would enable and start %s", service)
			continue
		}
		if _, err := helper.ExecSudo(nil, "", "systemctl", "enable", service); err != nil {
			return err
		}
		if _, err := helper.ExecSudo(nil, "", "systemctl", "start", service); err != nil {
			return err
		}
		out.Successf("Enabled %s", service)
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
			out.Infof("Would stop and disable %s", service)
			continue
		}
		if _, err := helper.ExecSudo(nil, "", "systemctl", "stop", service); err != nil {
			return err
		}
		if _, err := helper.ExecSudo(nil, "", "systemctl", "disable", service); err != nil {
			return err
		}
		out.Successf("Disabled %s", service)
	}
	return nil
}

// SystemdUserEnable enables a user systemd service
func SystemdUserEnable(in interface{}, out output.Output) error {
	out.InstructionTitle("Enabling user service")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, service := range data {
		if Dryrun {
			out.Infof("Would enable and start %s", service)
			continue
		}
		if _, err := helper.Exec(nil, "", "systemctl", "--user", "enable", service); err != nil {
			return err
		}
		if _, err := helper.Exec(nil, "", "systemctl", "--user", "start", service); err != nil {
			return err
		}
		out.Successf("Enabled %s", service)
	}
	return nil
}

// SystemdUserDisable disables a user systemd service
func SystemdUserDisable(in interface{}, out output.Output) error {
	out.InstructionTitle("Disabling user service")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, service := range data {
		if Dryrun {
			out.Infof("Would stop and disable %s", service)
			continue
		}
		if _, err := helper.Exec(nil, "", "systemctl", "--user", "stop", service); err != nil {
			return err
		}
		if _, err := helper.Exec(nil, "", "systemctl", "--user", "disable", service); err != nil {
			return err
		}
		out.Successf("Disabled %s", service)
	}
	return nil
}
