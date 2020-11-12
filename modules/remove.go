package modules

import (
	"os"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("remove", Remove)
}

func Remove(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Remove move file or (empty) directory")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				out.Info(path + " does not exist already")
				continue
			}
			return err
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		out.Success("Removed " + path)
	}
	return nil
}
