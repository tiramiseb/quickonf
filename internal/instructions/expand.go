package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
)

type Expand struct {
	Variable string
}

func (e *Expand) Name() string {
	return "expand"
}

func (e *Expand) Run(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	contentOfVar := vars.translateVariables("<" + e.Variable + ">")
	expanded := vars.translateVariables(contentOfVar)
	vars.define(e.Variable, expanded)
	return []*CheckReport{{
		"expand",
		commands.StatusSuccess,
		fmt.Sprintf("Expanded content of variable %s", e.Variable),
		nil,
		signalTarget,
		"",
		"",
	}}, true
}

func (e *Expand) Reset() {}
