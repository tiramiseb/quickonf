package conf

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

func (p *parser) command(toks []*token) instructions.Instruction {
	var targets []string
	for equalPos, tok := range toks {
		if tok.typ == tokenEqual {
			targets = make([]string, equalPos)
			for i := 0; i < equalPos; i++ {
				targets[i] = toks[i].content
			}
			toks = toks[equalPos+1:]
		}
	}
	commandName := toks[0].content
	args := make([]string, len(toks)-1)
	for i, tok := range toks[1:] {
		args[i] = tok.content
	}
	command, ok := commands.Get(commandName)
	if !ok {
		p.errs = append(p.errs, toks[0].errorf(`no command named "%s"`, commandName))
		return nil
	}
	if len(targets) > len(command.Outputs) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected maximum %d targets, got %d", len(command.Outputs), len(targets)),
		)
	}
	if len(args) != len(command.Arguments) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected %d arguments, got %d", len(command.Arguments), len(args)),
		)
		return nil
	}
	return &instructions.Command{Command: command, Arguments: args, Targets: targets}
}
