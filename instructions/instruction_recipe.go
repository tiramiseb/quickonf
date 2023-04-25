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
			message:      fmt.Sprintf("Recipe %q does not exist", r.RecipeName),
			signalTarget: signalTarget,
		}}, false
	}
	reports := []*CheckReport{{
		Name:         "recipe",
		level:        level,
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf("Running recipe %q...", r.RecipeName),
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

func (r *Recipe) NotRunReports(level int) []*CheckReport {
	rec, ok := recipes[r.RecipeName]
	if !ok {
		return []*CheckReport{{
			Name:    "recipe",
			level:   level,
			status:  commands.StatusNotRun,
			message: fmt.Sprintf("Recipe %q does not exist", r.RecipeName),
		}}
	}
	msg := r.description()
	reports := []*CheckReport{{
		Name:    "recipe",
		level:   level,
		status:  commands.StatusNotRun,
		message: msg.string(0),
	}}
	for _, ins := range rec.Instructions {
		thisReports := ins.NotRunReports(level + 1)
		reports = append(reports, thisReports...)
	}
	return reports
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
	content := r.description()
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

func (r *Recipe) description() stringBuilder {
	var content stringBuilder
	content.add("recipe")
	content.add(r.RecipeName)
	return content
}

func (r *Recipe) hasConfigError() bool {
	rec, ok := recipes[r.RecipeName]
	if !ok {
		return true
	}
	for _, ins := range rec.Instructions {
		if ins.hasConfigError() {
			return true
		}
	}
	return false
}
