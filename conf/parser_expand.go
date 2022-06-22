package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) expand(tokens tokens) (instrs []instructions.Instruction, next tokens) {
	if len(tokens) > 2 {
		p.errs = append(p.errs, tokens[0].error("expected a variable name as the only argument"))
		return nil, p.nextLine()
	}
	return []instructions.Instruction{&instructions.Expand{Variable: tokens[1].content}}, p.nextLine()
}
