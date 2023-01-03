package conf

import (
	"github.com/tiramiseb/quickonf/instructions"
)

type parser struct {
	groups    []*instructions.Group
	cookbooks []string

	tokens tokens
	idx    int

	errs []error
}

func newParser(tokens tokens) parser {
	return parser{
		tokens: tokens,
	}
}

func (p *parser) nextLine() (toks tokens) {
	for p.idx < len(p.tokens) && p.tokens[p.idx].typ != tokenEOL {
		toks = append(toks, p.tokens[p.idx])
		p.idx++
	}
	// Ignore tokenEOL
	p.idx++
	return toks
}

// parse parses the tokens in order to create a list of groups.
//
// All functions called from this function receive the "next" line for
// processing and return the "next" line for processing by another sub-parser.
// It is necessary in order to know how to process next line
func (p *parser) parse() (groups []*instructions.Group, err error) {
	next := p.nextLine()
	for next != nil {
		next = p.noIndentation(next)
	}
	if len(p.cookbooks) > 0 {
		maximumPriority := 0
		for _, g := range p.groups {
			if g.Priority > maximumPriority {
				maximumPriority = g.Priority
			}
		}
		cookbooks := &instructions.Group{
			Name:     "Cookbooks",
			Priority: maximumPriority + 1,
		}
		for _, uri := range p.cookbooks {
			cookbooks.Instructions = append(cookbooks.Instructions, &instructions.Cookbook{URI: uri, ReadFn: Read})
		}
		p.groups = append(p.groups, cookbooks)
	}
	return p.groups, nil
}
