package conf

import (
	"strconv"
)

func (c *checker) priority(toks tokens) (next tokens, newVars []string) {
	c.result.addToken(toks[0], CheckTypeKeyword)
	if len(toks) != 2 {
		c.result.addError(toks[0], CheckSeverityError, "Expected a priority value, as an integer")
	}
	_, err := strconv.Atoi(toks[1].content)
	if err != nil {
		c.result.addErrorf(toks[0], CheckSeverityError, "%s is not a valid integer", toks[1].content)
	}
	return c.nextLine(), nil
}
