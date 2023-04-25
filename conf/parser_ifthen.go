package conf

import (
	"errors"
	"path/filepath"

	"github.com/tiramiseb/quickonf/instructions"
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

func (p *parser) ifThen(toks tokens, group *instructions.Group, currentIndent int, knownVars map[string]string) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	p.checkResult.addToken(toks[0], CheckTypeKeyword)
	var operation instructions.Operation
	if len(toks) < 2 {
		instrs = append(instrs, toks[0].error("expected an operation"))
		p.checkResult.addError(toks[0], CheckSeverityError, "Expected an operation")
	} else {
		values := make([]tokenOrOperation, len(toks[1:]))
		for i, t := range toks[1:] {
			values[i] = tokenOrOperation{token: t}
		}
		values, errors := p.findConditionsWithOneArgument(values)
		if errors != nil {
			instrs = append(instrs, errors...)
		}
		values, errors = p.findConditionsComparison(values)
		if errors != nil {
			instrs = append(instrs, errors...)
		}
		if len(values) != 1 || values[0].operation == nil {
			instrs = append(instrs, values[0].token.error("Invalid condition"))
			p.checkResult.addError(values[0].token, CheckSeverityError, "Invalid condition")
		}
		operation = values[0].operation
	}
	next = p.nextLine()
	nextIndent, _ := next.indentation()
	if nextIndent <= currentIndent {
		instrs = append(instrs, toks[0].error(`expected commands in the if clause`))
		p.checkResult.addError(toks[0], CheckSeverityError, "Expected commands in the if clause")
	}
	inss, next, newVars := p.instructions(nil, next, group, nextIndent, knownVars)
	instrs = append(instrs, inss...)
	ins := &instructions.If{Operation: operation, Instructions: instrs}
	return []instructions.Instruction{ins}, next, newVars
}

func (p *parser) findConditionsWithOneArgument(input []tokenOrOperation) (output []tokenOrOperation, errorInstructions []instructions.Instruction) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsWithOneArgument[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case i >= len(input)-1:
			errorInstructions = append(errorInstructions, x.token.errorf("%s needs a file path", x.token.content))
			p.checkResult.addErrorf(x.token, CheckSeverityError, "%s needs a file path", x.token.content)
		default:
			i++
			y := input[i]
			switch {
			case y.operation != nil:
				errorInstructions = append(errorInstructions, y.token.errorf("%s needs a value", x.token.content))
				p.checkResult.addErrorf(y.token, CheckSeverityError, "%s needs a value", x.token.content)
			case !filepath.IsAbs(y.token.content):
				errorInstructions = append(errorInstructions, y.token.error("path must be absolute"))
				p.checkResult.addError(y.token, CheckSeverityError, "path must be absolute")
			default:
				op, err := cond(y.token.content)
				if err != nil {
					errorInstructions = append(errorInstructions, y.token.error(err.Error()))
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
	return output, errorInstructions
}

func (p *parser) findConditionsComparison(input []tokenOrOperation) (output []tokenOrOperation, errorInstructions []instructions.Instruction) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsComparison[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case output[len(output)-1].operation != nil:
			errorInstructions = append(errorInstructions, output[len(output)-1].token.errorf("must be a value (for %s)", x.token.content))
			p.checkResult.addErrorf(output[len(output)-1].token, CheckSeverityError, "must be a value (for %s)", x.token.content)
		default:
			left := output[len(output)-1].token
			i++
			if input[i].operation != nil {
				errorInstructions = append(errorInstructions, output[len(output)-1].token.errorf("must be a value (for %s)", x.token.content))
				p.checkResult.addErrorf(output[len(output)-1].token, CheckSeverityError, "must be a value (for %s)", x.token.content)
			}
			right := input[i].token
			op, err := cond(left.content, right.content)
			if err != nil {
				errorInstructions = append(errorInstructions, x.token.error(err.Error()))
			}
			output[len(output)-1] = tokenOrOperation{
				token:     x.token,
				operation: op,
			}
		}
		i++
	}
	return output, errorInstructions
}
