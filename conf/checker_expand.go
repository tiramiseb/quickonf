package conf

func (c *checker) expand(tokens tokens, knownVars []string) (next tokens, newVars []string) {
	c.result.addToken(tokens[0], CheckTypeKeyword)
	switch {
	case len(tokens) == 1:
		c.result.addError(tokens[0], CheckSeverityError, "expected a variable name as argument")
	case len(tokens) == 2:
		c.result.addVariableToken(tokens[1], knownVars)
	case len(tokens) > 2:
		c.result.addVariableToken(tokens[1], knownVars)
		for _, tok := range tokens[2:] {
			c.result.addError(tok, CheckSeverityWarning, "expected a variable name as the only argument")
		}
	}
	return c.nextLine(), nil
}
