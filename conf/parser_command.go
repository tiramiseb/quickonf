package conf

import (
	"strings"

	"github.com/tiramiseb/quickonf/commands"
	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) command(toks tokens, knownVars map[string]string) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	newVars = map[string]string{}
	var targets []string
	for equalPos, tok := range toks {
		if tok.typ == tokenEqual {
			assignmentInstructionSlice := make([]string, 0, len(toks)-equalPos-1)
			for i := equalPos + 1; i < len(toks); i++ {
				assignmentInstructionSlice = append(assignmentInstructionSlice, toks[i].content)
			}
			assignmentInstruction := strings.Join(assignmentInstructionSlice, " ")
			for i := 0; i < equalPos; i++ {
				p.checkResult.addToken(toks[i], CheckTypeVariable)
				newVars[toks[i].content] = assignmentInstruction
			}
			if equalPos == len(toks)-1 {
				p.checkResult.addError(tok, CheckSeverityError, "Missing command")
				return []instructions.Instruction{toks[equalPos].error("missing command")}, p.nextLine(), newVars
			}
			targets = make([]string, equalPos)
			for i := 0; i < equalPos; i++ {
				targets[i] = toks[i].content
			}
			toks = toks[equalPos+1:]
		}
	}
	commandToken := toks[0]
	commandName := commandToken.content
	args := make([]string, len(toks)-1)
	for i, tok := range toks[1:] {
		args[i] = tok.content
		p.checkResult.makeVarsTokens(tok, knownVars)
	}
	command, ok := commands.Get(commandName)
	if !ok {
		p.checkResult.addErrorf(commandToken, CheckSeverityError, `No command named "%s"`, commandName)
		p.checkResult.addUnknownCommandToken(toks[0])
		return []instructions.Instruction{toks[0].errorf(`no command named "%s"`, commandName)}, p.nextLine(), newVars
	}
	if command.LastOutputIsVariadic() {
		if len(targets) < len(command.Outputs) {
			instrs = append(
				instrs,
				toks[1].errorf("expected minimum %d targets, got %d", len(command.Outputs), len(targets)),
			)
			p.checkResult.addErrorf(commandToken, CheckSeverityError, "Expected minimum %d targets, got %d", len(command.Outputs), len(targets))
		}
	} else {
		if len(targets) != len(command.Outputs) {
			instrs = append(
				instrs,
				toks[1].errorf("expected %d targets, got %d", len(command.Outputs), len(targets)),
			)
			p.checkResult.addErrorf(commandToken, CheckSeverityError, "Expected %d targets, got %d", len(command.Outputs), len(targets))
		}
	}
	if command.LastArgumentIsVariadic() {
		if len(args) < len(command.Arguments) {
			instrs = append(
				instrs,
				toks[1].errorf("expected minimum %d arguments, got %d", len(command.Arguments), len(args)),
			)
			p.checkResult.addErrorf(commandToken, CheckSeverityError, "Expected minimum %d arguments, got %d", len(command.Arguments), len(args))
		}
	} else {
		if len(args) != len(command.Arguments) {
			instrs = append(
				instrs,
				toks[1].errorf("expected %d arguments, got %d", len(command.Arguments), len(args)),
			)
			p.checkResult.addErrorf(commandToken, CheckSeverityError, "Expected %d arguments, got %d", len(command.Arguments), len(args))
		}
	}
	p.checkResult.addCommandToken(toks[0], command)
	return append([]instructions.Instruction{&instructions.Command{Command: command, Arguments: args, Targets: targets}}, instrs...), p.nextLine(), newVars
}
