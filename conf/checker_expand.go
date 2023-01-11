package conf

func (c *checker) expand(tokens tokens, knownVars map[string]string) (next tokens, newVars map[string]string) {
	c.result.addToken(tokens[0], CheckTypeKeyword)
	switch {
	case len(tokens) == 1:
		c.result.addError(tokens[0], CheckSeverityError, "expected a variable name as argument")
	case len(tokens) == 2:
		c.addVarToken(tokens[1], knownVars)
	case len(tokens) > 2:
		c.addVarToken(tokens[1], knownVars)
		for _, tok := range tokens[2:] {
			c.result.addError(tok, CheckSeverityWarning, "expected a variable name as the only argument")
		}
	}
	return c.nextLine(), nil
}
