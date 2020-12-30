package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("directory", Directory)
	Register("root-directory", RootDirectory)
}

// Directory creates directories
func Directory(in interface{}, out output.Output) error {
	out.InstructionTitle("Directory creation")
	return directory(in, out, false)
}

// RootDirectory creates directories as root
func RootDirectory(in interface{}, out output.Output) error {
	out.InstructionTitle("Directory creation as root")
	return directory(in, out, true)

}

func directory(in interface{}, out output.Output, root bool) error {
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		result, err := helper.Directory(path, root)
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
