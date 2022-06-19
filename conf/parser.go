package conf

import (
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

type parser struct {
	groups  []*instructions.Group
	recipes map[string][]instructions.Instruction

	tokens tokens
	idx    int

	errs []error
}

func (p *parser) nextLine() (toks tokens) {
	for {
		if p.idx >= len(p.tokens) {
			return nil
		}
		t := p.tokens[p.idx]
		p.idx++
		if t.typ == tokenEOL {
			return toks
		}
		toks = append(toks, t)
	}
}

// parse parses the tokens in order to create a list of groups.
//
// All functions called from this function receive the "next" line for
// processing and return the "next" line for processing by another sub-parser.
// It is necessary in order to know how to process next line
func (p *parser) parse() ([]*instructions.Group, error) {
	next := p.nextLine()
	for next != nil {
		next = p.parseWithoutIndentation(next)
	}
	return p.groups, nil
}

func (p *parser) parseWithoutIndentation(line tokens) (next tokens) {
	next = p.nextLine()
	if len(line) == 0 {
		return
	}
	switch line[0].typ {
	case tokenGroupName:
		return p.parseGroup(line[0], next)
	case tokenCookbook:
		p.parseCookbook(line[0])
		return
	}

	// Illegal token...
	switch {
	case len(line) == 1:
		p.errs = append(p.errs, line[0].errorf(`expected group name, got "%s"`, line[0].content))
	case line[0].typ == tokenIndentation:
		p.errs = append(p.errs, line[1].errorf(`expected group name, got "%s"`, line[1].content))
	default:
		content := make([]string, len(line))
		for i, t := range line {
			content[i] = t.content
		}
		contentStr := strings.Join(content, " ")
		p.errs = append(p.errs, line[0].errorf(`expected group name, got "%s"`, contentStr))
	}
	return
}

func (p *parser) parseGroup(name *token, next tokens) tokens {
	group := &instructions.Group{Name: name.content}
	indent, _ := next.indentation()
	if indent == 0 {
		return next
	}
	group.Instructions, next = p.parseInstructions(nil, next, group, indent)
	nextIndent, firstToken := next.indentation()
	if nextIndent > 0 {
		p.errs = append(p.errs, firstToken.errorf("invalid indentation (expecting none or %d)", indent))
		next = p.nextLine()
	}
	p.groups = append(p.groups, group)
	return next
}

func (p *parser) parseInstructions(prefixAllWith []*token, line tokens, group *instructions.Group, currentIndent int) (instrs []instructions.Instruction, next tokens) {
	// Read a list of commands
	nbPrefixes := len(prefixAllWith)
	for {
		thisIndent, firstToken := line.indentation()
		if thisIndent > currentIndent {
			p.errs = append(p.errs, firstToken.errorf("invalid indentation (expecting %d)", currentIndent))
			next = p.nextLine()
			return
		} else if thisIndent < currentIndent {
			next = line
			return
		}
		if len(prefixAllWith) > 0 {
			newLine := make(tokens, len(prefixAllWith)+len(line))
			if line[0].typ == tokenIndentation {
				newLine[0] = line[0]
				for i, t := range prefixAllWith {
					newLine[i+1] = t
				}
				for i := 1; i < len(line); i++ {
					newLine[i+nbPrefixes] = line[i]
				}
			} else {
				copy(newLine, prefixAllWith)
				for i, t := range line {
					newLine[i+nbPrefixes] = t
				}
			}
			line = newLine
		}
		var ins instructions.Instruction
		switch firstToken.typ {
		case tokenExpand:
			if len(line) > 3 {
				p.errs = append(p.errs, firstToken.error("expected a variable name as the only argument"))
				break
			}
			ins = p.expand(line[2])
		case tokenIf:
			if len(line) < 3 {
				p.errs = append(p.errs, firstToken.error("expected an operation"))
				break
			}
			ins, next = p.ifThen(line[2:], group, currentIndent)
		case tokenPriority:
			p.priority(line[1:], group)
		case tokenRepeat:
			if len(line) < 3 {
				p.errs = append(p.errs, firstToken.error("expected something to repeat"))
				break
			}
			var inss []instructions.Instruction
			inss, next = p.repeat(line[2:], group, currentIndent)
			instrs = append(instrs, inss...)
		default:
			ins = p.command(line[1:])
		}
		if ins != nil {
			instrs = append(instrs, ins)
		}
		if next == nil {
			line = p.nextLine()
		} else {
			line = next
			next = nil
		}
	}
}
