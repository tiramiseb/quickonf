package modules

import (
	"path"

	quickonfErrors "github.com/tiramiseb/quickonf/internal/errors"
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("only-on", OnlyOn)
}

// OnlyOn runs the next instructions only if the current hostname matches the given string
func OnlyOn(in interface{}, out output.Output) error {
	out.InstructionTitle("Run only on the given host")
	data, err := helper.String(in)
	if err != nil {
		return err
	}
	host, err := helper.String("<hostname>")
	if err != nil {
		return err
	}
	ok, err := path.Match(data, host)
	if err != nil {
		return err
	}
	if ok {
		out.Successf(`Current hostname "%s" matches "%s"`, host, data)
		return nil
	}
	out.Infof(`Current hostname "%s" does not match "%s"`, host, data)
	return quickonfErrors.NoError
}
