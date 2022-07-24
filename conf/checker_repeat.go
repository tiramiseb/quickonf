package conf

func (c *checker) repeat(toks tokens, currentIndent int, knownVars []string) (next tokens, newVars []string) {
	if len(toks) < 2 {
		c.result.addError(toks[0], CheckSeverityError, "expected something to repeat")
		return c.nextLine(), nil
	}
	next = c.nextLine()
	indent, _ := next.indentation()
	if indent <= currentIndent {
		c.result.addError(toks[0], CheckSeverityError, "expected arguments in the repeat clause")
		return next, nil
	}
	return c.instructions(toks[1:], next, indent, knownVars)
}
