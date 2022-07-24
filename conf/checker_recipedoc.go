package conf

func (c *checker) recipeDoc(toks tokens) (next tokens, newVars []string) {
	if len(toks) != 2 {
		c.result.addError(toks[0], CheckSeverityError, "expected a documentation")
	}
	return c.nextLine(), nil
}

func (c *checker) recipeVarDoc(toks tokens) (next tokens, newVars []string) {
	if len(toks) >= 2 {
		c.result.addToken(toks[1], CheckTypeVariable)
	}
	if len(toks) != 3 {
		c.result.addError(toks[0], CheckSeverityError, "expected a variable name and a documentation")
	}
	return c.nextLine(), nil
}
