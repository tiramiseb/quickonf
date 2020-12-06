package modules

import (
	"errors"
	"os"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("symlink", Symlink)
}

// Symlink creates symbolic links
func Symlink(in interface{}, out output.Output) error {
	out.InstructionTitle("Create symbolic links")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for path, target := range data {
		path = helper.Path(path)
		target = helper.Path(target)
		if stat, err := os.Lstat(target); err == nil {
			if stat.Mode()&os.ModeSymlink != os.ModeSymlink {
				return errors.New(path + " already exists but is not a symlink")
			}
			if currentTarget, err := os.Readlink(path); err == nil {
				if currentTarget == target {
					out.Info("Link from " + path + " to " + target + " already exists")
					return nil
				}
				if Dryrun {
					out.Info("Would remove " + path)
				} else {
					if err := os.Remove(path); err != nil {
						return err
					}
				}
			} else if !os.IsNotExist(err) {
				return err
			}

		} else if !os.IsNotExist(err) {
			return err
		}
		if Dryrun {
			out.Info("Would create link from " + path + " to " + target)
			continue
		}
		if err := os.Symlink(target, path); err != nil {
			return err
		}
		out.Success("Link from " + path + " to " + target + " created")
	}
	return nil
}
