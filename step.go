package quickonf

import (
	"errors"

	"github.com/tiramiseb/quickonf/modules"
	"github.com/tiramiseb/quickonf/output"
	quickonfErrors "github.com/tiramiseb/quickonf/errors"
)

type message struct {
	isError bool
	message string
}

type Step map[string][]map[string]interface{}

func (step Step) Run(out output.Output) {
	store := map[string]interface{}{}
	for title, actions := range step {
		out.StepTitle(title)
		for _, action := range actions {
			if err := RunAction(action, out, store); err != nil {
                if err != quickonfErrors.NoError {
    				out.Error(err)
                }
				return
			}
		}
		return
	}
}

func RunAction(action map[string]interface{}, out output.Output, store map[string]interface{}) error {
	for name, data := range action {
		module := modules.Get(name)
		if module == nil {
			return errors.New("[No module named \"" + name + "\"]")
		}
		return module(data, out, store)
	}
	return errors.New("[No action]")
}
