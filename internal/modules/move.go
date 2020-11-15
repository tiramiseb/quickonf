package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("move", Move)
	Register("force-move", ForceMove)
}

// Move moves files or directories, or does nothing if the source does not exist
func Move(in interface{}, out output.Output) error {
	out.InstructionTitle("Move file or directory")
	return move(in, out, false)
}

// ForceMove moves files or directories, removing the destination if it exists
func ForceMove(in interface{}, out output.Output) error {
	out.InstructionTitle("Move file or directory, crushing destination if necessary")
	return move(in, out, true)
}

func move(in interface{}, out output.Output, force bool) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for from, to := range data {
		from = helper.Path(from)
		to = helper.Path(to)
		if _, err := os.Stat(from); err != nil {
			if os.IsNotExist(err) {
				out.Info("Source " + from + " does not exist")
				continue
			}
			return err
		}
		_, err := os.Stat(to)
		if err == nil {
			if !force {
				return errors.New(to + " already exists")
			}
			if err := os.RemoveAll(to); err != nil {
				return err
			}
		} else {
			if !os.IsNotExist(err) {
				return err
			}
		}
		if err = os.Rename(from, to); err != nil {
			return err
		}
		out.Success("Moved " + from + " to " + to)
	}
	return nil
}
