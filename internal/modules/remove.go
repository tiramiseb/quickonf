package modules

import (
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("remove", Remove)
}

// Remove removes files or (empty) directories
func Remove(in interface{}, out output.Output) error {
	out.InstructionTitle("Remove file or (empty) directory")
	data, err := helper.SliceString(in)
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
