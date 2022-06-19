package conf

import (
	"strconv"

	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) priority(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens) {
	if len(toks) != 2 {
		p.errs = append(p.errs, toks[0].error("expected a priority value, as an integer"))
	}
	priority, err := strconv.Atoi(toks[1].content)
	if err != nil {
		p.errs = append(p.errs, toks[1].errorf("%s is not a valid integer", toks[1]))
	}
	group.Priority = priority
	return nil, p.nextLine()
}
