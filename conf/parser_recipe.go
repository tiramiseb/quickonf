package conf

import (
	"strings"

	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) recipe(toks tokens, group *instructions.Group, currentIndent int) ([]instructions.Instruction, tokens, map[string]string) {
	recipeInstruction := &instructions.Recipe{}
	otherInstrs := []instructions.Instruction{}
	if len(toks) >= 2 {
		content := make([]string, len(toks[1:]))
		for i, t := range toks[1:] {
			content[i] = t.content
		}
		recipeInstruction.RecipeName = strings.Join(content, " ")
	} else {
		otherInstrs = append(otherInstrs, toks[0].error("expected a recipe name"))
		p.checkResult.addError(toks[0], CheckSeverityError, "expected a recipe name")
	}

	vars := map[string]string{}
	next := p.nextLine()
	varIndent, varTokens := next.indentation()
	firstVarIndent := varIndent
	for varIndent > currentIndent {
		if varIndent != firstVarIndent {
			otherInstrs = append(otherInstrs, varTokens[0].error("wrong indentation"))
		}
		switch {
		case len(varTokens) == 0:
			otherInstrs = append(otherInstrs, next[0].error("expected a variable assignment (var = value)"))
		case len(varTokens) != 3:
			otherInstrs = append(otherInstrs, varTokens[0].error("expected a variable assignment (var = value)"))
		case varTokens[1].typ != tokenEqual:
			otherInstrs = append(otherInstrs, varTokens[1].error("expected a variable assignment (var = value)"))
		default:
			vars[varTokens[0].content] = varTokens[2].content
		}
		next = p.nextLine()
		varIndent, varTokens = next.indentation()
	}
	recipeInstruction.Variables = vars
	return append([]instructions.Instruction{recipeInstruction}, otherInstrs...), next, nil
}
