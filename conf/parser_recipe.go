package conf

import (
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) recipe(toks tokens, group *instructions.Group, currentIndent int) ([]instructions.Instruction, tokens) {
	if len(toks) < 2 {
		p.errs = append(p.errs, toks[0].error("expected a recipe name"))
		return nil, p.nextLine()
	}

	content := make([]string, len(toks[1:]))
	for i, t := range toks[1:] {
		content[i] = t.content
	}
	recipeName := strings.Join(content, " ")

	vars := map[string]string{}
	next := p.nextLine()
	varIndent, varTokens := next.indentation()
	firstVarIndent := varIndent
	for varIndent > currentIndent {
		if varIndent != firstVarIndent {
			p.errs = append(p.errs, varTokens[0].error("wrong indentation"))
		}
		switch {
		case len(varTokens) == 0:
			p.errs = append(p.errs, next[0].error("expected a variable assignment (var = value)"))
		case len(varTokens) != 3:
			p.errs = append(p.errs, varTokens[0].error("expected a variable assignment (var = value)"))
		case varTokens[1].typ != tokenEqual:
			p.errs = append(p.errs, varTokens[1].error("expected a variable assignment (var = value)"))
		default:
			vars[varTokens[0].content] = varTokens[2].content
		}
		next = p.nextLine()
		varIndent, varTokens = next.indentation()
	}
	return []instructions.Instruction{&instructions.Recipe{RecipeName: recipeName, Variables: vars}}, next
}
