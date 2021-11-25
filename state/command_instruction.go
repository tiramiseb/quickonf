package state

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/output"
)

type Instruction struct {
	Instruction instructions.Instruction
	Arguments   []string
	Targets     []string
}

func (i *Instruction) Run(groupOut *output.Group, vars variables, options Options) bool {
	out := groupOut.NewInstruction(i.Instruction.Name)
	out.Info("Running...")
	args := make([]string, len(i.Arguments))
	for i, src := range i.Arguments {
		args[i] = vars.translateVariables(src)
	}
	slow(options)
	result, ok := i.Instruction.Run(args, out)
	for i, tgt := range i.Targets {
		if len(result) <= i {
			break
		}
		vars.define(tgt, result[i])
	}
	return ok
}
