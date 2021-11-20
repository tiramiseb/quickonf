package conf

import (
	"io"

	"github.com/tiramiseb/quickonf/state"
)

func Read(r io.Reader, filters []string) (*state.State, []error) {
	lxr := newLexer(r)

	tokens, err := lxr.scan()
	if err != nil {
		return nil, []error{err}
	}
	p := parser{tokens: tokens}
	groups := p.parse()
	return &state.State{Groups: groups}, p.errs
}
