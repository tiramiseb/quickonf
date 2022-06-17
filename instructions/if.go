package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/commands"
)

type If struct {
	Operation    Operation
	Instructions []Instruction
}

func (i *If) Name() string {
	return "if"
}

func (i *If) RunCheck(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	success := i.Operation.Compare(vars)
	if !success {
		return []*CheckReport{{
			Name:         "if",
			status:       commands.StatusSuccess,
			message:      fmt.Sprintf(`"%s" is false, not running contained instructions...`, i.Operation.String()),
			signalTarget: signalTarget,
		}}, true
	}
	reports := []*CheckReport{{
		Name:         "if",
		status:       commands.StatusSuccess,
		message:      fmt.Sprintf(`"%s" is true, running contained instructions...`, i.Operation.String()),
		signalTarget: signalTarget,
	}}
	for _, ins := range i.Instructions {
		thisReports, ok := ins.RunCheck(vars, signalTarget)
		if thisReports != nil {
			reports = append(reports, thisReports...)
		}
		if !ok {
			return reports, false
		}
	}
	return reports, true
}

func (i *If) Reset() {
	for _, ins := range i.Instructions {
		ins.Reset()
	}
}
