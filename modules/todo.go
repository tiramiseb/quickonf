package modules

import (
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("todo", Todo)
}

func Todo(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("TODO")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	for _, msg := range data {
		out.Alert(msg)
	}
	return nil
}
