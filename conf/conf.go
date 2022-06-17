package conf

import (
	"io"

	"github.com/tiramiseb/quickonf/instructions"
)

func Read(r io.Reader) (*instructions.Groups, []error) {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return nil, []error{err}
	}
	p := parser{tokens: tokens}
	groups, err := p.parse()
	if err != nil {
		return nil, []error{err}
	}
	return instructions.NewGroups(groups), p.errs
}
