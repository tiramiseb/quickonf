package conf

import (
	"io"

	"github.com/tiramiseb/quickonf/internal/state"
)

func Read(r io.Reader, filters []string) (*state.State, []error) {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return nil, []error{err}
	}
	p := parser{tokens: tokens}
	groups, err := p.parse(filters)
	if err != nil {
		return nil, []error{err}
	}
	return &state.State{Filtered: len(filters) > 0, Groups: groups}, p.errs
}
