package main

import (
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

func (e *embeder) command(level int, cmd *instructions.Command) {
	e.write(level, "&Command{")
	e.write(level+1, "Command: commands.UGet(\"%s\"),", cmd.Command.Name)
	if len(cmd.Arguments) > 0 {
		e.write(level+1, "Arguments: []string{")
		for _, a := range cmd.Arguments {
			e.write(level+2, "`%s`,", strings.ReplaceAll(a, "`", "\\`"))
		}
		e.write(level+1, "},")
	}
	if len(cmd.Targets) > 0 {
		e.write(level+1, "Targets: []string{")
		for _, t := range cmd.Targets {
			e.write(level+2, "`%s`,", strings.ReplaceAll(t, "`", "\\`"))
		}
		e.write(level+1, "},")
	}
	e.write(level, "},")
}
