package state

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/output"
)

type Argument struct {
	IsVariable bool   // True to read data from the variable
	Value      string // Content of the argument or variable name
}

// TODO More dynamic stuff:
type Command struct {
	Group       *Group
	Name        string
	Instruction instructions.Instruction
	Arguments   []Argument
	Targets     []string
}

// Run executes the instruction for this command and returns true if it succeeds
func (c *Command) Run(groupOut *output.Group) bool {
	out := groupOut.NewInstruction(c.Name)
	out.Info("Running...")
	result, ok := c.Instruction(c.argsToStrings(), out)
	for i, t := range c.Targets {
		c.Group.setVariable(t, result[i])
	}
	return ok

}

func (c *Command) argsToStrings() []string {
	strs := make([]string, len(c.Arguments))
	for i, arg := range c.Arguments {
		if arg.IsVariable {
			strs[i] = c.Group.getVariable(arg.Value)
		} else {
			strs[i] = arg.Value
		}
	}
	return strs
}
