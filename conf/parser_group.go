package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) group(name *token) (next tokens) {
	p.checkResult.addToken(name, CheckTypeGroupname)
	firstInstruction := p.nextLine()
	indent, _ := firstInstruction.indentation()
	if indent == 0 {
		// The group was empty, this line has no indentation, ignore the group and process that line
		p.checkResult.addError(name, CheckSeverityInformation, "This group is empty")
		return firstInstruction
	}

	group := &instructions.Group{Name: name.content}
	group.Instructions, next, _ = p.instructions(nil, firstInstruction, group, indent, instructions.GlobalVars())
	p.groups = append(p.groups, group)

	nextIndent, remainingTokens := next.indentation()
	if nextIndent > 0 {
		p.groups = append(p.groups, &instructions.Group{
			Name: "No group",
			Instructions: []instructions.Instruction{
				remainingTokens[0].errorf("invalid indentation (expecting none or %d)", indent),
			},
		})
		p.checkResult.addErrorf(remainingTokens[0], CheckSeverityError, "invalid indentation (expecting none or %d)", indent)
		next = p.nextLine()
	}
	return next
}
