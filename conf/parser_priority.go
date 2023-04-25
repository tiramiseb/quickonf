package conf

import (
	"strconv"

	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) priority(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	if len(toks) != 2 {
		instrs = append(instrs, toks[0].error("expected a priority value, as an integer"))
		p.checkResult.addError(toks[0], CheckSeverityError, "Expected a priority value, as an integer")
	}
	priority, err := strconv.Atoi(toks[1].content)
	if err != nil {
		instrs = append(instrs, toks[1].errorf("%s is not a valid integer", toks[1]))
		p.checkResult.addErrorf(toks[0], CheckSeverityError, "%s is not a valid integer", toks[1].content)
	}
	group.Priority = priority
	return nil, p.nextLine(), nil
}
