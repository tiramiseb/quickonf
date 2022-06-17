package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) repeat(toks []*token, group *instructions.Group, currentIndent int) ([]instructions.Instruction, tokens) {
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, toks[0].error(`expected arguments in the repeat clause`))
	}
	return p.parseInstructions(toks, next, group, indent)
}
