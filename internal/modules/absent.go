package modules

import (
	"errors"
	"io"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("absent", Absent)
}

// Absent makes sure a file is absent
func Absent(in interface{}, out output.Output) error {
	out.InstructionTitle("Make file absent")
	data, err := helper.SliceString(in)
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
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.Readdirnames(1)
			if err == nil {
				return errors.New("Directory " + path + " contains files, cannot be deleted")
			}
			if err != io.EOF {
				return err
			}
		}
		if Dryrun {
			out.Info("Would remove " + path)
			continue
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}
		out.Success(path + " removed")
	}
	return nil
}
