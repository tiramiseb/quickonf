package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands"
)

type Expand struct {
	Variable string
}

func (e *Expand) Name() string {
	return "expand"
}

func (e *Expand) RunCheck(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	contentOfVar := vars.TranslateVariables("<" + e.Variable + ">")
	expanded := vars.TranslateVariables(contentOfVar)
	vars.define(e.Variable, expanded)
	return []*CheckReport{{
		Name:         "expand",
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf("Expanded content of variable %s", e.Variable),
		signalTarget: signalTarget,
	}}, true
}

func (e *Expand) Reset() {}
