package modules

import (
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("store", Store)
}

// Store adds values in the store
func Store(in interface{}, out output.Output) error {
	out.InstructionTitle("Store")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	for k, v := range data {
		helper.Store(k, v)
		out.Infof("Stored %s in %s", v, k)
	}
	return nil
}
