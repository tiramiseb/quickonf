package conf

import (
	"strconv"

	"github.com/tiramiseb/quickonf/internal/instructions"
)

func (p *parser) priority(toks []*token, group *instructions.Group) {
	if len(toks) != 2 {
		p.errs = append(p.errs, toks[0].error("expected a priority value, as an integer"))
	}
	priority, err := strconv.Atoi(toks[1].content)
	if err != nil {
		p.errs = append(p.errs, toks[1].errorf("%s is not a valid integer", toks[1]))
	}
	group.Priority = priority
}
