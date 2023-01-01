package conf

func (c *checker) ifThen(toks tokens, currentIndent int, knownVars map[string]string) (next tokens, newVars map[string]string) {
	c.result.addToken(toks[0], CheckTypeKeyword)
	if len(toks) < 2 {
		c.result.addError(toks[0], CheckSeverityError, "Expected an operation")
	} else {
		values := make([]tokenOrOperation, len(toks[1:]))
		for i, t := range toks[1:] {
			values[i] = tokenOrOperation{token: t}
		}
		values = c.findConditionsWithOneArgument(values, knownVars)
		values = c.findConditionsComparison(values, knownVars)
		if len(values) != 1 || values[0].operation == nil {
			c.result.addError(values[0].token, CheckSeverityError, "Invalid condition")
		}
	}
	next = c.nextLine()
	nextIndent, _ := next.indentation()
	if nextIndent <= currentIndent {
		c.result.addError(toks[0], CheckSeverityError, "Expected commands in the if clause")
	}
	return c.instructions(nil, next, nextIndent, knownVars)
}

func (c *checker) findConditionsWithOneArgument(input []tokenOrOperation, knownVars map[string]string) (output []tokenOrOperation) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsWithOneArgument[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case i >= len(input)-1:
			c.result.addErrorf(x.token, CheckSeverityError, "%s needs a file path", x.token.content)
		default:
			i++
			y := input[i]
			switch {
			case y.operation != nil:
				c.result.addErrorf(y.token, CheckSeverityError, "%s needs a value", x.token.content)
			default:
				c.result.addToken(x.token, CheckTypeFunction)
				c.makeVarsTokens(y.token, knownVars)
				op, err := cond(y.token.content)
				if err != nil {
					c.result.addError(y.token, CheckSeverityError, err.Error())
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

func (c *checker) findConditionsComparison(input []tokenOrOperation, knownVars map[string]string) (output []tokenOrOperation) {
	i := 0
	for i < len(input) {
		x := input[i]
		cond, ok := conditionsComparison[x.token.content]
		switch {
		case x.operation != nil, !ok:
			output = append(output, x)
		case output[len(output)-1].operation != nil:
			c.result.addErrorf(output[len(output)-1].token, CheckSeverityError, "must be a value (for %s)", x.token.content)
		default:
			left := output[len(output)-1].token
			i++
			if input[i].operation != nil {
				c.result.addErrorf(output[len(output)-1].token, CheckSeverityError, "must be a value (for %s)", x.token.content)
			}
			right := input[i].token
			c.result.addToken(x.token, CheckTypeOperator)
			c.makeVarsTokens(left, knownVars)
			c.makeVarsTokens(right, knownVars)
			op, err := cond(left.content, right.content)
			if err != nil {
				c.result.addError(x.token, CheckSeverityError, err.Error())
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
