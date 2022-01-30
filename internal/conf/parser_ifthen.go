package conf

import (
	"errors"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/instructions"
)

type tokenOrOperation struct {
	token     *token
	operation instructions.Operation
}

var conditionsWithOneArgument = map[string]func(arg string) (instructions.Operation, error){
	"file.absent": func(arg string) (instructions.Operation, error) {
		if !filepath.IsAbs(arg) {
			return nil, errors.New("path must be absolute")
		}
		return &instructions.FileAbsent{Path: arg}, nil
	},
	"file.present": func(arg string) (instructions.Operation, error) {
		if !filepath.IsAbs(arg) {
			return nil, errors.New("path must be absolute")
		}
		return &instructions.FilePresent{Path: arg}, nil
	},
}

var conditionsComparison = map[string]func(left, right string) (instructions.Operation, error){
	"=": func(left, right string) (instructions.Operation, error) {
		return &instructions.Equal{Left: left, Right: right}, nil
	},
	"!=": func(left, right string) (instructions.Operation, error) {
		return &instructions.Different{Left: left, Right: right}, nil
	},
}

func (p *parser) ifThen(toks []*token, group *instructions.Group, currentIndent int) (instructions.Instruction, tokens) {
	values := make([]tokenOrOperation, len(toks))
	for i, t := range toks {
		values[i] = tokenOrOperation{token: t}
	}
	values = p.findConditionsWithOneArgument(values)
	values = p.findConditionsComparison(values)
	if len(values) != 1 || values[0].operation == nil {
		p.errs = append(p.errs, values[0].token.error("Invalid condition"))
	}
	next := p.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		p.errs = append(p.errs, toks[0].error(`expected commands in the if clause`))
	}
	inss, next := p.parseInstructions(nil, next, group, indent)
	ins := &instructions.If{Operation: values[0].operation, Instructions: inss}
	return ins, next
}

func (p *parser) findConditionsWithOneArgument(input []tokenOrOperation) (output []tokenOrOperation) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsWithOneArgument[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case i >= len(input)-1:
			p.errs = append(p.errs, x.token.errorf("%s needs a file path", x.token.content))
		default:
			i++
			y := input[i]
			switch {
			case y.operation != nil:
				p.errs = append(p.errs, y.token.errorf("%s needs a value", x.token.content))
			case !filepath.IsAbs(y.token.content):
				p.errs = append(p.errs, y.token.error("path must be absolute"))
			default:
				op, err := cond(y.token.content)
				if err != nil {
					p.errs = append(p.errs, y.token.error(err.Error()))
				}
				output = append(output,
					tokenOrOperation{
						token:     x.token,
						operation: op,
					},
				)
			}
		}
		i++
	}
	return output
}

func (p *parser) findConditionsComparison(input []tokenOrOperation) (output []tokenOrOperation) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsComparison[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case output[len(output)-1].operation != nil:
			p.errs = append(p.errs, output[len(output)-1].token.errorf("must be a value (for %s)", x.token.content))
		default:
			left := output[len(output)-1].token
			i++
			if input[i].operation != nil {
				p.errs = append(p.errs, output[len(output)-1].token.errorf("must be a value (for %s)", x.token.content))
			}
			right := input[i].token
			op, err := cond(left.content, right.content)
			if err != nil {
				p.errs = append(p.errs, x.token.error(err.Error()))
			}
			output[len(output)-1] = tokenOrOperation{
				token:     x.token,
				operation: op,
			}
		}
		i++
	}
	return output
}
