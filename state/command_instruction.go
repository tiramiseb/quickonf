package state

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/output"
)

type Instruction struct {
	Instruction instructions.Instruction
	Arguments   []string
}

func (i *Instruction) Run(groupOut *output.Group, vars variables) bool {
	out := groupOut.NewInstruction(i.Instruction.Name)
	out.Info("Running...")
	args := make([]string, len(i.Arguments))
	for i, src := range i.Arguments {
		args[i] = vars.translateVariables(src)
	}
	_, ok := i.Instruction.Run(args, out)
	return ok
}
