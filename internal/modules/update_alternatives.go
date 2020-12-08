package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("update-alternatives", UpdateAlternatives)
}

// UpdateAlternatives changes default commands
func UpdateAlternatives(in interface{}, out output.Output) error {
	out.InstructionTitle("Update Alternative")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for alt, path := range data {
		if Dryrun {
			out.Info("Would change alternative for " + alt + " to " + path)
			continue
		}
		if _, err := helper.ExecSudo("update-alternatives", "--set", alt, path); err != nil {
			return err
		}
		out.Success("Changed alternative for " + alt + " to " + path)
	}
	return nil
}
