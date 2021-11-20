package conf

import (
	"github.com/tiramiseb/quickonf/state"
)

type parser struct {
	tokens []*token
	idx    int

	errs []error
}

func (p *parser) nextLine() (tokens []*token) {
	tokens = []*token{}
	for {
		if p.idx >= len(p.tokens) {
			return nil
		}
		token := p.tokens[p.idx]
		p.idx++
		if token.typ == tokenEOL {
			return tokens
		}
		tokens = append(tokens, token)
	}
}

func (p *parser) parse() (groups []*state.Group) {
	// for _, t := range p.tokens {
	// 	DDcontent := fmt.Sprintf("%s", t.content)
	// 	if len(DDcontent) > 70 {
	// 		DDcontent = DDcontent[:70]
	// 	}
	// 	fmt.Printf("[%d:%d] %d >> %s\n", t.line, t.column, t.typ, DDcontent)
	// }
	next := p.nextLine()
	for {
		if next == nil {
			break
		}
		var group *state.Group
		group, next = p.parseGroup(next)
		if group != nil {
			groups = append(groups, group)
		}
	}
	return
}

func (p *parser) parseGroup(line []*token) (group *state.Group, next []*token) {
	next = p.nextLine()
	if len(line) == 0 {
		return
	}
	if line[0].typ != tokenGroupName {
		// Expecting group name only. In any other situation, it is a syntax error
		switch {
		case len(line) == 1:
			p.errs = append(p.errs, line[0].errorf(`expected group name, got "%v"`, line[0].content))
		case line[0].typ == tokenIndentation:
			p.errs = append(p.errs, line[1].errorf(`expected group name, got "%v"`, line[1].content))
		default:
			p.errs = append(p.errs, line[0].error("expected group name, got an empty line"))
		}
		return
	}
	group = &state.Group{Name: line[0].content.(string)}
	// Next line MUST be an indentation
	if len(next) == 0 || next[0].typ != tokenIndentation {
		return
	}
	group.Commands = p.parseCommands(next, next[0].content.(int))
	return
}

func (p *parser) parseCommands(line []*token, indent int) []*state.Command {
	// Read a list of commands
	return nil
}
