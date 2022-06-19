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

func newParser(tokens tokens) parser {
	return parser{
		tokens:  tokens,
		recipes: map[string][]instructions.Instruction{},
	}
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
	if len(line) == 0 {
		// Line is empty, ignore it (should not happen, empty lines are removed by the lexer)
		return p.nextLine()
	}
	switch line[0].typ {
	case tokenGroupName:
		return p.parseGroup(line[0], p.nextLine())
	case tokenCookbook:
		p.parseCookbook(line[0])
		return p.nextLine()
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
	return p.nextLine()
}

func (p *parser) parseGroup(name *token, firstInstruction tokens) (next tokens) {
	indent, _ := firstInstruction.indentation()
	if indent == 0 {
		// The group was empty, this line has no indentation, process it as a new group
		return firstInstruction
	}

	group := &instructions.Group{Name: name.content}
	group.Instructions, next = p.parseInstructions(nil, firstInstruction, group, indent)
	p.groups = append(p.groups, group)

	nextIndent, remainingTokens := next.indentation()
	if nextIndent > 0 {
		p.errs = append(p.errs, remainingTokens[0].errorf("invalid indentation (expecting none or %d)", indent))
		next = p.nextLine()
	}
	return next
}

func (p *parser) parseInstructions(prefixAllWith, line tokens, group *instructions.Group, currentIndent int) (instrs []instructions.Instruction, next tokens) {
	// Read a list of instructions...
	for {
		if len(line) == 0 {
			// End of file (the lexer doesn't leave any empty line)
			return
		}
		thisIndent, toks := line.indentation()

		switch {
		case thisIndent > currentIndent:
			// A larger indentation should not happen, unless we are in another block (which would be started when parsing an instruction needing another block)
			p.errs = append(p.errs, toks[0].errorf("invalid indentation (expecting %d)", currentIndent))
			return instrs, p.nextLine()
		case thisIndent < currentIndent:
			// This indentation block is finished, quit the function
			return instrs, line
		}

		// Add prefix to instruction if needed
		toks = addPrefix(prefixAllWith, toks)

		// Parse the tokens for this instruction!
		var ins []instructions.Instruction
		switch toks[0].typ {
		case tokenExpand:
			ins, line = p.expand(toks)
		case tokenIf:
			ins, line = p.ifThen(toks, group, currentIndent)
		case tokenPriority:
			ins, line = p.priority(toks, group)
		case tokenRepeat:
			ins, line = p.repeat(toks, group, currentIndent)
		default:
			ins, line = p.command(toks)
		}
		instrs = append(instrs, ins...)
	}
}

func addPrefix(prefix tokens, existing tokens) tokens {
	if len(prefix) == 0 {
		return existing
	}

	nbPrefixes := len(prefix)

	newTokens := make(tokens, len(prefix)+len(existing))

	copy(newTokens, prefix)

	for i, t := range existing {
		newTokens[i+nbPrefixes] = t
	}

	return newTokens
}
