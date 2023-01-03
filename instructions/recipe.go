package instructions

import (
	"fmt"
	"strings"

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

func (r *Recipe) RunCheck(vars *Variables, signalTarget chan bool, level int) ([]*CheckReport, bool) {
	rec, ok := recipes[r.RecipeName]
	if !ok {
		return []*CheckReport{{
			Name:         "recipe",
			level:        level,
			status:       commands.StatusError,
			message:      fmt.Sprintf(`"Recipe "%s"`, r.RecipeName),
			signalTarget: signalTarget,
		}}, false
	}
	reports := []*CheckReport{{
		Name:         "recipe",
		level:        level,
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf(`Running recipe "%s"...`, r.RecipeName),
		signalTarget: signalTarget,
	}}

	r.instructions = make([]Instruction, len(rec.Instructions))
	copy(r.instructions, rec.Instructions)

	thisVars := vars.clone()
	for key, value := range r.Variables {
		value = vars.TranslateVariables(value)
		thisVars.define(key, value)
	}

	for _, ins := range r.instructions {
		thisReports, ok := ins.RunCheck(thisVars, signalTarget, level+1)
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

func (r *Recipe) String() string {
	return r.indentedString(0)
}

func (r *Recipe) indentedString(level int) string {
	var result []string
	var content stringBuilder
	content.add("recipe")
	content.add(r.RecipeName)
	result = append(result, content.string(level))
	for key, value := range r.Variables {
		var content stringBuilder
		content.add(key)
		content.add("=")
		content.add(value)
		result = append(result, content.string(level+1))
	}
	return strings.Join(result, "\n")
}
