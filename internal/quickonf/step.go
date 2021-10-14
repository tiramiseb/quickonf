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

// Name returns the step name
func (step Step) Name() string {
	for title := range step {
		return title
	}
	return ""
}

// Always returns true if the step must always run
func (step Step) Always() bool {
	for _, instructions := range step {
		for _, instruction := range instructions {
			for k, v := range instruction {
				if k == "always" {
					b, ok := v.(bool)
					if ok {
						return b
					}
					s, ok := v.(string)
					if ok {
						return s == "true"
					}
				}

			}

		}
	}
	return false
}

func (step Step) run(out output.Output) {
	var skipNext bool
	for title, instructions := range step {
		out.StepTitle(title)
	instruction:
		for _, instruction := range instructions {
			if skipNext {
				skipNext = false
				continue
			}
			for k := range instruction {
				if k == "always" {
					continue instruction
				}
			}
			if err := runAction(instruction, out); err != nil {
				if err == quickonfErrors.NoError {
					return
				}
				if err == quickonfErrors.SkipNext {
					skipNext = true
					continue
				}
				out.Error(err)
			}
		}
		return
	}
}

func (step Step) list(out output.Output) {
	for title := range step {
		out.Info(title)
		return
	}
}

func runAction(action map[string]interface{}, out output.Output) error {
	for instructionName, data := range action {
		return modules.Get(instructionName)(data, out)
	}
	return errors.New("[No instruction]")
}
