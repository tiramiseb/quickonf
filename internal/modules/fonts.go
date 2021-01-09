package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("font-cache", FontCache)
}

// FontCache regenerates the font cache
func FontCache(in interface{}, out output.Output) error {
	out.InstructionTitle("Regenerate font cache")
	if Dryrun {
		out.Info("Would regenerate the font cache")
	}
	out.ShowLoader()
	_, err := helper.Exec(nil, "", "fc-cache", "-f")
	out.HideLoader()
	if err != nil {
		return err
	}
	out.Info("Font cache regenerated")
	return nil
}
