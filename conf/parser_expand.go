package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) expand(tokens tokens, knownVars map[string]string) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	p.checkResult.addToken(tokens[0], CheckTypeKeyword)
	switch {
	case len(tokens) == 1:
		instrs = []instructions.Instruction{tokens[0].error("expected a variable name as argument")}
		p.checkResult.addError(tokens[0], CheckSeverityError, "expected a variable name as argument")
	case len(tokens) == 2:
		instrs = []instructions.Instruction{&instructions.Expand{Variable: tokens[1].content}}
		p.checkResult.addVariableToken(tokens[1], knownVars)
	case len(tokens) > 2:
		instrs = []instructions.Instruction{tokens[0].error("expected a variable name as the only argument")}
		p.checkResult.addVariableToken(tokens[1], knownVars)
		for _, tok := range tokens[2:] {
			p.checkResult.addError(tok, CheckSeverityWarning, "expected a variable name as the only argument")
		}
	}
	return instrs, p.nextLine(), nil
}
