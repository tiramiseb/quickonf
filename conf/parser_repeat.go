package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) repeat(toks tokens, group *instructions.Group, currentIndent int) ([]instructions.Instruction, tokens) {
	if len(toks) < 2 {
		p.errs = append(p.errs, toks[0].error("expected something to repeat"))
		return nil, p.nextLine()
	}
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, toks[0].error(`expected arguments in the repeat clause`))
		return nil, next
	}
	return p.parseInstructions(toks[1:], next, group, indent)
}
