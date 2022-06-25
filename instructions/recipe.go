package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands"
)

type Recipe struct {
	RecipeName string
	Variables  map[string]string

	instructions []Instruction
}

func (r *Recipe) Name() string {
	return "recipe"
}

func (r *Recipe) RunCheck(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	rec, ok := recipes[r.RecipeName]
	if !ok {
		return []*CheckReport{{
			Name:         "recipe",
			status:       commands.StatusError,
			message:      fmt.Sprintf(`"Recipe "%s"`, r.RecipeName),
			signalTarget: signalTarget,
		}}, false
	}
	reports := []*CheckReport{{
		Name:         "recipe",
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf(`Running recipe "%s"...`, r.RecipeName),
		signalTarget: signalTarget,
	}}

	r.instructions = make([]Instruction, len(rec))
	copy(r.instructions, rec)

	thisVars := vars.clone()
	for key, value := range r.Variables {
		value = vars.TranslateVariables(value)
		thisVars.define(key, value)
	}

	for _, ins := range r.instructions {
		thisReports, ok := ins.RunCheck(thisVars, signalTarget)
		if thisReports != nil {
			reports = append(reports, thisReports...)
		}
		if !ok {
			return reports, false
		}
	}
	return reports, true

}

func (r *Recipe) Reset() {
	for _, ins := range r.instructions {
		ins.Reset()
	}
}
