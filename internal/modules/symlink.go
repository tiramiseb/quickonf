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
			out.Infof("Link from %s to %s already exists", path, target)
		case helper.SymlinkCreated:
			if Dryrun {
				out.Infof("Would create link from %s to %s", path, target)
			} else {
				out.Successf("Link from %s to %s created", path, target)
			}
		}
	}
	return nil
}
