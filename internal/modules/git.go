package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("git-config", GitConfig)
}

// GitConfig sets got configuration parameters
func GitConfig(in interface{}, out output.Output) error {
	out.InstructionTitle("Git configuration")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for param, value := range data {
		if Dryrun {
			out.Info("Would set " + param + " to " + value)
			continue
		}
		if _, err := helper.Exec(nil, "git", "config", "--global", param, value); err != nil {
			return err
		}
		out.Success("Set " + param + " to " + value)
	}
	return nil
}
