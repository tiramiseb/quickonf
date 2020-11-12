package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("directory", Directory)
}

func Directory(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Directory creation")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		info, err := os.Lstat(path)
		if err == nil {
			if info.IsDir() {
				out.Info(path + " already exists")
				continue
			}
			return errors.New(path + " is not a directory")
		}
		if !os.IsNotExist(err) {
			return err
		}
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		out.Success(path + " created")
	}
	return nil
}
