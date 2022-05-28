package conf

import (
	"io"
	"sort"

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
	// Sort groups by priority
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Priority > groups[j].Priority
	})
	return groups, p.errs
}
