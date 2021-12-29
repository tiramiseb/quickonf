package commands

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands/datastores"
	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(alternativesSet)
}

var alternativesSet = Command{
	"alternatives.set",
	"Set an alternative",
	[]string{
		"Name of the alternative",
		"Value to set (file path)",
	},
	nil,
	"Use Vim\n  alternatives.set vi /usr/bin/vim.basic",
	func(args []string) (result []string, msg string, apply *Apply, status Status) {
		name := args[0]
		value := args[1]

		current, err := datastores.Alternatives.Get(name)
		if err != nil {
			return nil, fmt.Sprintf("Could not get current value for %s: %s", name, err), nil, StatusError
		}
		if current == value {
			return nil, fmt.Sprintf("Current value for %s is already %s", name, current), nil, StatusSuccess
		}
		apply = &Apply{
			"alternatives.set",
			fmt.Sprintf("Will set alternative %s to %s", name, value),
			func(out Output) bool {
				out.Infof("Setting alternative %s to %s", name, value)
				wait, err := helper.Exec(nil, nil, "update-alternatives", "--set", name, value)
				if err != nil {
					out.Errorf("Could not change alternative: %s", err)
					return false
				}
				if err := wait(); err != nil {
					out.Errorf("Could not change alternative: %s", err)
					return false
				}
				out.Successf("Set alternative %s to %s", name, err)
				return true
			},
		}
		return nil, fmt.Sprintf("Need to set alternative %s to %s", name, value), apply, StatusInfo
	},
	datastores.Alternatives.Reset,
}
