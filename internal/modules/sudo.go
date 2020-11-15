package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("sudo-password", SudoPassword)
}

// SudoPassword sets the password for sudo
func SudoPassword(in interface{}, out output.Output) error {
	out.InstructionTitle("Set password for sudo")
	data, err := helper.String(in)
	if err != nil {
		return err
	}
	helper.SudoPassword = data
	out.Success("Password stored in memory (not displayed here, of course)")
	return nil
}
