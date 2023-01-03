package conf

import "strings"

func (p *parser) noIndentation(line tokens) (next tokens) {
	if len(line) == 0 {
		// Line is empty, ignore it (should not happen, empty lines are removed by the lexer)
		return p.nextLine()
	}
	switch line[0].typ {
	case tokenGroupName:
		return p.group(line[0])
	case tokenCookbook:
		return p.cookbook(line)
	}

	// Illegal token...
	switch {
	case len(line) == 1:
		p.errs = append(p.errs, line[0].errorf(`expected group name, got "%s"`, line[0].content))
	case line[0].typ == tokenIndentation:
		p.errs = append(p.errs, line[1].errorf(`expected group name, got "%s"`, line[1].content))
	default:
		content := make([]string, len(line))
		for i, t := range line {
			content[i] = t.content
		}
		contentStr := strings.Join(content, " ")
		p.errs = append(p.errs, line[0].errorf(`expected group name, got "%s"`, contentStr))
	}
	return p.nextLine()
}
