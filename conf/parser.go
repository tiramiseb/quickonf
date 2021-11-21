package conf

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/state"
)

type parser struct {
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

func (p *parser) parseGroup(line tokens) (group *state.Group, next tokens) {
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
	group = &state.Group{Name: line[0].content}
	// Next line MUST start with an indentation
	indent, _ := next.indentation()
	if indent == 0 {
		return
	}
	group.Commands, next = p.parseCommands(next, indent)
	nextIndent, firstToken := next.indentation()
	if nextIndent > 0 {
		p.errs = append(p.errs, firstToken.errorf("invalid indentation (expecting none or %d)", indent))
		next = p.nextLine()
	}
	return
}

func (p *parser) parseCommands(line tokens, currentIndent int) (commands []state.Command, next tokens) {
	// Read a list of commands
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
		var cmd state.Command
		switch firstToken.typ {
		case tokenIf:
			cmd, next = p.ifThen(line[2:], currentIndent)
		default:
			cmd = p.instruction(line[1:])
		}
		if cmd != nil {
			commands = append(commands, cmd)
		}
		if next == nil {
			line = p.nextLine()
		} else {
			line = next
			next = nil
		}
	}
}

func (p *parser) ifThen(toks []*token, currentIndent int) (state.Command, tokens) {
	// Later, add support for "and", "or", etc
	left := toks[0]
	operator := toks[1]
	right := toks[2]
	if left.typ != tokenDefault {
		p.errs = append(p.errs, left.errorf(`expected value, got "%s"`, left.content))
	}
	if !operator.isOperator() {
		p.errs = append(p.errs, operator.errorf(`expected operator, got "%s"`, operator.content))
	}
	if right.typ != tokenDefault {
		p.errs = append(p.errs, right.errorf(`expected value, got "%s"`, right.content))
	}
	var operation state.Operation
	switch operator.typ {
	case tokenEqual:
		operation = &state.Equal{Left: left.content, Right: right.content}
	}
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, operator.error(`expected commands in the if clause`))
	}
	cmds, next := p.parseCommands(next, indent)
	cmd := &state.If{Operation: operation, Commands: cmds}
	return cmd, next
}

func (p *parser) instruction(toks []*token) state.Command {
	var targets []string
	for equalPos, tok := range toks {
		if tok.typ == tokenEqual {
			for i := 0; i < equalPos; i++ {
				targets[i] = toks[i].content
			}
			toks = toks[equalPos+1:]
		}
	}
	instructionName := toks[0].content
	args := make([]string, len(toks)-1)
	for i, tok := range toks[1:] {
		args[i] = tok.content
	}
	instruction, ok := instructions.Get(instructionName)
	if !ok {
		p.errs = append(p.errs, toks[0].errorf(`no instruction named "%s"`, instructionName))
		return nil
	}
	if len(targets) > len(instruction.Outputs) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected maximum %d targets, got %d", len(instruction.Outputs), len(targets)),
		)
	}
	if len(args) != len(instruction.Arguments) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected %d arguments, got %d", len(instruction.Arguments), len(args)),
		)
		return nil
	}
	return &state.Instruction{Instruction: instruction, Arguments: args, Targets: targets}
}
