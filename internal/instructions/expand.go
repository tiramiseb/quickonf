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

func (e *Expand) Run(vars Variables) ([]commands.Apply, []CheckReport, bool) {
	contentOfVar := vars.translateVariables("<" + e.Variable + ">")
	expanded := vars.translateVariables(contentOfVar)
	vars.define(e.Variable, expanded)
	return nil, []CheckReport{{
		Name:    "expand",
		Status:  commands.StatusInfo,
		Message: fmt.Sprintf("Expanded content of variable %s", e.Variable),
	}}, true
}

func (e *Expand) Reset() {}
