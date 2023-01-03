package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) group(name *token) (next tokens) {
	firstInstruction := p.nextLine()
	indent, _ := firstInstruction.indentation()
	if indent == 0 {
		// The group was empty, this line has no indentation, ignore the group and process that line
		return firstInstruction
	}

	group := &instructions.Group{Name: name.content}
	group.Instructions, next = p.instructions(nil, firstInstruction, group, indent)
	p.groups = append(p.groups, group)

	nextIndent, remainingTokens := next.indentation()
	if nextIndent > 0 {
		p.errs = append(p.errs, remainingTokens[0].errorf("invalid indentation (expecting none or %d)", indent))
		next = p.nextLine()
	}
	return next
}
