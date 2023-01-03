package conf

import (
	"fmt"
	"io"

	"github.com/tiramiseb/quickonf/instructions"
)

func Read(r io.Reader) (*instructions.Groups, []error) {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return nil, []error{err}
	}
	p := newParser(tokens)
	groups, err := p.parse()
	if err != nil {
		return nil, []error{err}
	}
	return instructions.NewGroups(groups), p.errs
}

func Check(r io.Reader) error {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return err
	}
	c := newChecker(tokens)
	result := c.check()
	json, err := result.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(json)
	return nil
}
