package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("absent", Absent)
}

func Absent(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Make file absent")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, path := range data {
		path = helper.Path(path)
		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				out.Info(path + " already absent")
				continue
			}
			return err
		}
		if info.IsDir() {
			ok, err := helper.IsEmpty(path)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("Directory " + path + " contains files, cannot be deleted")
			}
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}
		out.Success(path + " removed")
	}
	return nil
}
