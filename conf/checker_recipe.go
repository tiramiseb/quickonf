package conf

func (c *checker) recipe(toks tokens, currentIndent int, knownVars []string) (next tokens, newVars []string) {
	c.result.addToken(toks[0], CheckTypeKeyword)
	if len(toks) < 2 {
		c.result.addError(toks[0], CheckSeverityError, "expected a recipe name")
	}

	content := make([]string, len(toks[1:]))
	for i, t := range toks[1:] {
		content[i] = t.content
	}

	next = c.nextLine()
	varIndent, varTokens := next.indentation()
	firstVarIndent := varIndent
	for varIndent > currentIndent {
		if varIndent != firstVarIndent {
			c.result.addError(varTokens[0], CheckSeverityWarning, "Wrong indentation")
		}
		switch {
		case len(varTokens) == 0:
			c.result.addError(next[0], CheckSeverityError, "Expected a variable assignment (var = value)")
		case len(varTokens) != 3:
			c.result.addError(varTokens[0], CheckSeverityError, "Expected a variable assignment (var = value)")
		case varTokens[1].typ != tokenEqual:
			c.result.addError(varTokens[1], CheckSeverityError, "Expected a variable assignment (var = value)")
		default:
			c.result.addToken(varTokens[0], CheckTypeVariable)
			c.makeVarsTokens(varTokens[2], knownVars)
		}
		next = c.nextLine()
		varIndent, varTokens = next.indentation()
	}
	return next, nil
}
