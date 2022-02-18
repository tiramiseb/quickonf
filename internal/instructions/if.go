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

func (i *If) Run(vars Variables) ([]CheckReport, bool) {
	success := i.Operation.Compare(vars)
	if !success {
		return []CheckReport{{
			"if",
			commands.StatusInfo,
			fmt.Sprintf(`"%s" is false, not running contained instructions...`, i.Operation.String()),
			nil,
		}}, true
	}
	reports := []CheckReport{
		{
			"if",
			commands.StatusInfo,
			fmt.Sprintf(`"%s" is true, running contained instructions...`, i.Operation.String()),
			nil,
		},
	}
	for _, ins := range i.Instructions {
		thisReports, ok := ins.Run(vars)
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
