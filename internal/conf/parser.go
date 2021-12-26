package conf

import (
	"strconv"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
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
func (p *parser) parse() (groups []*instructions.Group, err error) {
	next := p.nextLine()
	for {
		if next == nil {
			break
		}
		var group *instructions.Group
		group, next = p.parseGroup(next)
		if group != nil {
			groups = append(groups, group)
		}
	}
	return
}

func (p *parser) parseGroup(line tokens) (group *instructions.Group, next tokens) {
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

	group = &instructions.Group{Name: line[0].content}
	indent, _ := next.indentation()
	if indent == 0 {
		return
	}
	group.Instructions, next = p.parseInstructions(nil, next, group, indent)
	nextIndent, firstToken := next.indentation()
	if nextIndent > 0 {
		p.errs = append(p.errs, firstToken.errorf("invalid indentation (expecting none or %d)", indent))
		next = p.nextLine()
	}
	return
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
				for i, t := range prefixAllWith {
					newLine[i] = t
				}
				for i, t := range line {
					newLine[i+nbPrefixes] = t
				}
			}
			line = newLine
		}
		var ins instructions.Instruction
		switch firstToken.typ {
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
		case tokenIf:
			if len(line) < 3 {
				p.errs = append(p.errs, firstToken.error("expected an operation"))
				break
			}
			ins, next = p.ifThen(line[2:], group, currentIndent)
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

func (p *parser) priority(toks []*token, group *instructions.Group) {
	if len(toks) != 2 {
		p.errs = append(p.errs, toks[0].error("expected a priority value, as an integer"))
	}
	priority, err := strconv.Atoi(toks[1].content)
	if err != nil {
		p.errs = append(p.errs, toks[1].errorf("%s is not a valid integer", toks[1]))
	}
	group.Priority = priority
}

func (p *parser) ifThen(toks []*token, group *instructions.Group, currentIndent int) (instructions.Instruction, tokens) {
	// Later, add support for "and", "or", etc
	if len(toks) != 3 {
		p.errs = append(p.errs, toks[0].error("expected a value followed by an operator followed by another value"))
		return nil, nil
	}
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
	var operation instructions.Operation
	switch operator.typ {
	case tokenEqual:
		operation = &instructions.Equal{Left: left.content, Right: right.content}
	}
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, operator.error(`expected commands in the if clause`))
	}
	inss, next := p.parseInstructions(nil, next, group, indent)
	ins := &instructions.If{Operation: operation, Instructions: inss}
	return ins, next
}

func (p *parser) repeat(toks []*token, group *instructions.Group, currentIndent int) ([]instructions.Instruction, tokens) {
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, toks[0].error(`expected arguments in the repeat clause`))
	}
	return p.parseInstructions(toks, next, group, indent)
}

func (p *parser) command(toks []*token) instructions.Instruction {
	var targets []string
	for equalPos, tok := range toks {
		if tok.typ == tokenEqual {
			for i := 0; i < equalPos; i++ {
				targets[i] = toks[i].content
			}
			toks = toks[equalPos+1:]
		}
	}
	commandName := toks[0].content
	args := make([]string, len(toks)-1)
	for i, tok := range toks[1:] {
		args[i] = tok.content
	}
	command, ok := commands.Get(commandName)
	if !ok {
		p.errs = append(p.errs, toks[0].errorf(`no command named "%s"`, commandName))
		return nil
	}
	if len(targets) > len(command.Outputs) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected maximum %d targets, got %d", len(command.Outputs), len(targets)),
		)
	}
	if len(args) != len(command.Arguments) {
		p.errs = append(
			p.errs,
			toks[1].errorf("expected %d arguments, got %d", len(command.Arguments), len(args)),
		)
		return nil
	}
	return &instructions.Command{Command: command, Arguments: args, Targets: targets}
}
