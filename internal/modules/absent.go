package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("absent", Absent)
}

// Absent makes sure a file is absent
func Absent(in interface{}, out output.Output) error {
	out.InstructionTitle("Make sure file is absent")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		status, err := helper.Remove(path)
		switch status {
		case helper.ResultError:
			return err
		case helper.ResultAlready:
			out.Infof("%s already absent", path)
		case helper.ResultDryrun:
			out.Infof("Would remove %s", path)
		case helper.ResultSuccess:
			out.Successf("%s removed", path)
		}
	}
	return nil
}
