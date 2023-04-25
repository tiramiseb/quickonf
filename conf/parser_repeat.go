package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) repeat(toks tokens, group *instructions.Group, currentIndent int, knownVars map[string]string) ([]instructions.Instruction, tokens, map[string]string) {
	if len(toks) < 2 {
		p.checkResult.addError(toks[0], CheckSeverityError, "expected something to repeat")
		return []instructions.Instruction{toks[0].error("expected something to repeat")}, p.nextLine(), nil
	}
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.checkResult.addError(toks[0], CheckSeverityError, "expected arguments in the repeat clause")
		return []instructions.Instruction{toks[0].error("expected arguments in the repeat clause")}, next, nil
	}
	return p.instructions(toks[1:], next, group, indent, knownVars)
}
