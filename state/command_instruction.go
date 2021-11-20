package state

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/output"
)

type Instruction struct {
	Instruction instructions.Instruction
	Arguments   []string
}

func (i *Instruction) Run(groupOut *output.Group) bool {
	out := groupOut.NewInstruction(i.Instruction.Name)
	out.Info("Running...")
	_, ok := i.Instruction.Run(i.Arguments, out)
	return ok
}
