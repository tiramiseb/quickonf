package conf

import (
	"github.com/tiramiseb/quickonf/instructions"
)

func (c *checker) group(name *token) (next tokens) {
	c.result.addToken(name, CheckTypeGroupname)
	firstInstruction := c.nextLine()
	indent, _ := firstInstruction.indentation()
	if indent == 0 {
		// The group was empty, this line has no indentation, ignore the group and process that line
		c.result.addError(name, CheckSeverityInformation, "This group is empty")
		return firstInstruction
	}

	next, _ = c.instructions(nil, firstInstruction, indent, instructions.GlobalVars())

	nextIndent, remainingTokens := next.indentation()
	if nextIndent > 0 {
		c.result.addErrorf(remainingTokens[0], CheckSeverityError, "invalid indentation (expecting none or %d)", indent)
		next = c.nextLine()
	}
	return next
}
