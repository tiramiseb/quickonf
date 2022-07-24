package conf

import "strings"

func (c *checker) noIndentation(line tokens) (next tokens) {
	if len(line) == 0 {
		// Line is empty, ignore it (should not happen, empty lines are removed by the lexer)
		return c.nextLine()
	}
	switch line[0].typ {
	case tokenGroupName:
		return c.group(line[0])
	case tokenCookbook:
		return c.cookbook(line)
	}

	// Illegal token...
	switch {
	case len(line) == 1:
		c.result.addErrorf(line[0], CheckSeverityError, `expected group name, got "%s"`, line[0].content)
	case line[0].typ == tokenIndentation:
		c.result.addErrorf(line[1], CheckSeverityError, `expected group name, got "%s"`, line[1].content)
	default:
		content := make([]string, len(line))
		for i, t := range line {
			content[i] = t.content
		}
		contentStr := strings.Join(content, " ")
		mixToken := token{
			line:    line[0].line,
			column:  line[0].column,
			length:  line[len(line)-1].column + line[len(line)-1].column - 1 - line[0].column,
			typ:     tokenDefault,
			content: contentStr,
		}
		c.result.addErrorf(&mixToken, CheckSeverityError, `expected group name, got "%s"`, contentStr)
	}
	return c.nextLine()
}
