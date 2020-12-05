package quickonf

import (
	"errors"

	quickonfErrors "github.com/tiramiseb/quickonf/internal/errors"
	"github.com/tiramiseb/quickonf/internal/modules"
	"github.com/tiramiseb/quickonf/internal/output"
)

type message struct {
	isError bool
	message string
}

// Step is a step definition, which includes instructions
type Step map[string][]map[string]interface{}

func (step Step) run(out output.Output, mode string) {
	for title, instructions := range step {
		switch mode {
		case "title":
			out.Info(title)
		case "action":
			out.StepTitle(title)
			for _, instruction := range instructions {
				if err := runAction(instruction, out); err != nil {
					if err != quickonfErrors.NoError {
						out.Error(err)
					}
					return
				}
			}
		}
		return
	}
}

func runAction(action map[string]interface{}, out output.Output) error {
	for name, data := range action {
		instruction := modules.Get(name)
		if instruction == nil {
			return errors.New("[No instruction named \"" + name + "\"]")
		}
		return instruction(data, out)
	}
	return errors.New("[No instruction]")
}
