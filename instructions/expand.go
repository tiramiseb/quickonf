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

func (e *Expand) RunCheck(vars Variables, signalTarget chan bool, level int) ([]*CheckReport, bool) {
	contentOfVar := vars.TranslateVariables("<" + e.Variable + ">")
	expanded := vars.TranslateVariables(contentOfVar)
	vars.define(e.Variable, expanded)
	return []*CheckReport{{
		Name:         "expand",
		level:        level,
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf("Expanded content of variable %s", e.Variable),
		signalTarget: signalTarget,
	}}, true
}

func (e *Expand) Reset() {}

func (e *Expand) String() string {
	return e.indentedString(0)
}

func (e *Expand) indentedString(level int) string {
	var content stringBuilder
	content.add("expand")
	content.add(e.Variable)
	return content.string(level)

}
