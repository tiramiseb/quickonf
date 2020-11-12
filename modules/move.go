package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/modules/helper"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("move", Move)
	Register("force-move", ForceMove)
}

func Move(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Move file or directory")
	return move(in, out, false, store)
}

func ForceMove(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Move file or directory, crushing destination if necessary")
	return move(in, out, true, store)
}

func move(in interface{}, out output.Output, force bool, store map[string]interface{}) error {
	data, err := input.MapStringString(in, store)
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
			if err := os.Remove(to); err != nil {
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
