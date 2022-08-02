package conf

import (
	"github.com/tiramiseb/quickonf/commands"
)

func (c *checker) command(toks tokens, knownVars []string) (next tokens, newVars []string) {
	var targetsCount int
	for pos, tok := range toks {
		if tok.typ == tokenEqual {
			targetsCount = pos
			for i := 0; i < pos; i++ {
				c.result.addToken(toks[i], CheckTypeVariable)
				newVars = append(newVars, toks[i].content)
			}
			if pos == len(toks)-1 {
				c.result.addError(tok, CheckSeverityError, "Missing command")
				return c.nextLine(), newVars
			}
			toks = toks[pos+1:]
			break
		}
	}
	commandToken := toks[0]
	commandName := commandToken.content
	argsCount := len(toks[1:])
	for _, tok := range toks[1:] {
		c.makeVarsTokens(tok, knownVars)
	}
	command, ok := commands.Get(commandName)
	if !ok {
		c.result.addErrorf(commandToken, CheckSeverityError, `No command named "%s"`, commandName)
		c.result.addUnknownCommandToken(toks[0])
		return c.nextLine(), nil
	}
	if targetsCount > len(command.Outputs) {
		c.result.addErrorf(commandToken, CheckSeverityError, "Expected maximum %d targets, got %d", len(command.Outputs), targetsCount)
	}
	if argsCount != len(command.Arguments) {
		c.result.addErrorf(commandToken, CheckSeverityError, "Expected %d arguments, got %d", len(command.Arguments), argsCount)
	}
	c.result.addCommandToken(toks[0], command)
	return c.nextLine(), newVars
}
