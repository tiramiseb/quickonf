package modules

import (
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
		status, err := helper.Symlink(path, target)
		switch status {
		case helper.SymlinkError:
			return err
		case helper.SymlinkAleradyExists:
			out.Info("Link from " + path + " to " + target + " already exists")
		case helper.SymlinkCreated:
			if Dryrun {
				out.Info("Would create link from " + path + " to " + target)
			} else {
				out.Success("Link from " + path + " to " + target + " created")
			}
		}
	}
	return nil
}
