package instructions

import (
	"fmt"
	"strings"

	"github.com/tiramiseb/quickonf/commands"
)

type If struct {
	Operation    Operation
	Instructions []Instruction
}

func (i *If) Name() string {
	return "if"
}

func (i *If) RunCheck(vars *Variables, signalTarget chan bool, level int) ([]*CheckReport, bool) {
	success := i.Operation.Compare(vars)
	if !success {
		return []*CheckReport{{
			Name:         "if",
			level:        level,
			status:       commands.StatusSuccess,
			message:      fmt.Sprintf(`"%s" is false, not running contained instructions...`, i.Operation.String()),
			signalTarget: signalTarget,
		}}, true
	}
	reports := []*CheckReport{{
		Name:         "if",
		level:        level,
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf(`"%s" is true, running contained instructions...`, i.Operation.String()),
		signalTarget: signalTarget,
	}}
	for _, ins := range i.Instructions {
		thisReports, ok := ins.RunCheck(vars, signalTarget, level+1)
		reports = append(reports, thisReports...)
		if !ok {
			return reports, false
		}
	}
	return reports, true
}

func (i *If) NotRunReports(level int) []*CheckReport {
	msg := i.description()
	reports := []*CheckReport{{
		Name:    "if",
		level:   level,
		status:  commands.StatusNotRun,
		message: msg.string(0),
	}}
	for _, ins := range i.Instructions {
		thisReports := ins.NotRunReports(level + 1)
		reports = append(reports, thisReports...)
	}
	return reports
}

func (i *If) Reset() {
	for _, ins := range i.Instructions {
		ins.Reset()
	}
}

func (i *If) String() string {
	return i.indentedString(0)
}

func (i *If) indentedString(level int) string {
	var result []string
	content := i.description()
	result = append(result, content.string(level))
	for _, instr := range i.Instructions {
		result = append(result, instr.indentedString(level+1))
	}
	return strings.Join(result, "\n")
}

func (i *If) description() stringBuilder {
	var content stringBuilder
	content.add("if")
	content.add(i.Operation.String())
	return content
}

func (i *If) hasConfigError() bool {
	for _, ins := range i.Instructions {
		if ins.hasConfigError() {
			return true
		}
	}
	return false
}
