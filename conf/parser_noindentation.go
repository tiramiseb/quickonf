package conf

import (
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

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
		p.groups = append(p.groups, &instructions.Group{
			Name: "No group",
			Instructions: []instructions.Instruction{
				line[0].errorf(`expected group name, got "%s"`, line[0].content),
			},
		})
		p.checkResult.addErrorf(line[0], CheckSeverityError, `expected group name, got "%s"`, line[0].content)
	case line[0].typ == tokenIndentation:
		p.groups = append(p.groups, &instructions.Group{
			Name: "No group",
			Instructions: []instructions.Instruction{
				line[1].errorf(`expected group name, got "%s"`, line[1].content),
			},
		})
		p.checkResult.addErrorf(line[1], CheckSeverityError, `expected group name, got "%s"`, line[1].content)
	default:
		content := make([]string, len(line))
		for i, t := range line {
			content[i] = t.content
		}
		contentStr := strings.Join(content, " ")
		p.groups = append(p.groups, &instructions.Group{
			Name: "No group",
			Instructions: []instructions.Instruction{
				line[1].errorf(`expected group name, got "%s"`, contentStr),
			},
		})
		p.checkResult.addErrorf(
			&token{
				line:    line[0].line,
				column:  line[0].column,
				length:  line[len(line)-1].column + line[len(line)-1].column - 1 - line[0].column,
				typ:     tokenDefault,
				content: contentStr,
			},
			CheckSeverityError,
			`expected group name, got "%s"`, contentStr,
		)
	}
	return p.nextLine()
}
