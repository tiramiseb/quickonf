package modules

import (
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("make", Make)
}

// Make executes the make command (using a makefile) in the given directories, with the given arguments
func Make(in interface{}, out output.Output) error {
	out.InstructionTitle("Make with a Makefile")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for cd, args := range data {
		if Dryrun {
			out.Infof("Would make \"%s\" in %s", args, cd)
			continue
		}
		out.ShowLoader()
		argz := []string{"-C", cd}
		argz = append(argz, strings.Fields(args)...)
		_, err := helper.Exec(nil, "", "make", argz...)
		out.HideLoader()
		if err != nil {
			return err
		}
		out.Successf("Made \"%s\" in %s", args, cd)
	}
	return nil
}
