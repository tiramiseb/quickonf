package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("directory", Directory)
}

// Directory creates directories
func Directory(in interface{}, out output.Output) error {
	out.InstructionTitle("Directory creation")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		result, err := helper.Directory(path)
		switch result {
		case helper.ResultAlready:
			out.Infof("%s already exists", path)
		case helper.ResultDryrun:
			out.Infof("Would create %s", path)
		case helper.ResultError:
			return err
		case helper.ResultSuccess:
			out.Successf("%s created", path)
		}
	}
	return nil
}
