package modules

import (
	"errors"
	"regexp"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("regexp-replace", RegexpReplace)
}

// RegexpReplace replaces
func RegexpReplace(in interface{}, out output.Output) error {
	out.InstructionTitle("Regex")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	from, ok := data["from"]
	if !ok {
		return errors.New("Missing from")
	}
	reg, ok := data["regexp"]
	if !ok {
		return errors.New("Missing regexp")
	}
	repl, ok := data["replace"]
	if !ok {
		return errors.New("Missing replace")
	}

	re, err := regexp.Compile(reg)
	if err != nil {
		return err
	}

	result := re.ReplaceAllString(from, repl)
	out.Info("Tranforming " + from)
	out.Info("Result is " + result)
	store, ok := data["store"]
	if ok {
		helper.Store(store, result)
	}

	return nil
}
