package instructions

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
)

type If struct {
	Operation    Operation
	Instructions []Instruction
}

func (i *If) Name() string {
	return "if"
}

func (i *If) Run(vars Variables) ([]commands.Apply, []CheckReport, bool) {
	success := i.Operation.Compare(vars)
	if !success {
		return nil, []CheckReport{{
			"if",
			commands.StatusInfo,
			fmt.Sprintf(`"%s" is false, not running contained instructions...`, i.Operation.String()),
		}}, true
	}
	var applies []commands.Apply
	reports := []CheckReport{
		{
			"if",
			commands.StatusInfo,
			fmt.Sprintf(`"%s" is true, running contained instructions...`, i.Operation.String()),
		},
	}
	for _, ins := range i.Instructions {
		thisApplies, thisReports, ok := ins.Run(vars)
		if thisApplies != nil {
			applies = append(applies, thisApplies...)
		}
		if thisReports != nil {
			reports = append(reports, thisReports...)
		}
		if !ok {
			return applies, reports, false
		}
	}
	return applies, reports, true
}

func (i *If) Reset() {
	for _, ins := range i.Instructions {
		ins.Reset()
	}
}
