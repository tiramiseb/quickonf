package modules

import (
	"fmt"
	"os"

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
		info, err := os.Lstat(path)
		if err == nil {
			if info.IsDir() {
				out.Infof("%s already exists", path)
				continue
			}
			return fmt.Errorf("%s is not a directory", path)
		}
		if !os.IsNotExist(err) {
			return err
		}
		if Dryrun {
			out.Infof("Would create %s", path)
			continue
		}
		if root {
			if _, err = helper.ExecSudo(nil, "mkdir", "-p", path); err != nil {
				return err
			}

		} else {
			if err = os.MkdirAll(path, 0755); err != nil {
				return err
			}
		}
		out.Successf("%s created", path)
	}
	return nil
}
