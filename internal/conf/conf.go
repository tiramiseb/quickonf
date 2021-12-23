package conf

import (
	"io"

	"github.com/tiramiseb/quickonf/internal/instructions"
)

func Read(r io.Reader) ([]*instructions.Group, []error) {
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
	return groups, p.errs
}
