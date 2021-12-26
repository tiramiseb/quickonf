package conf

import "github.com/tiramiseb/quickonf/internal/instructions"

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
