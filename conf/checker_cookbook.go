package conf

func (c *checker) cookbook(line tokens) (next tokens) {
	c.result.addToken(line[0], CheckTypeKeyword)
	if len(line) == 1 {
		c.result.addError(line[0], CheckSeverityError, "Cookbook URI not provided")
	}
	return c.nextLine()

}
