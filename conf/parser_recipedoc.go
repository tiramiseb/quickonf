package conf

import (
	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) recipeDoc(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	if len(toks) != 2 {
		instrs = append(instrs, toks[0].error("expected a documentation"))
		p.checkResult.addError(toks[0], CheckSeverityError, "expected a documentation")
	}
	group.RecipeDoc = toks[1].content
	return instrs, p.nextLine(), nil
}

func (p *parser) recipeVarDoc(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	if len(toks) >= 2 {
		p.checkResult.addToken(toks[1], CheckTypeVariable)
	}
	if len(toks) != 3 {
		instrs = append(instrs, toks[0].error("expected a variable name and a documentation"))
		p.checkResult.addError(toks[0], CheckSeverityError, "expected a variable name and a documentation")
	}
	if group.RecipeVarsDoc == nil {
		group.RecipeVarsDoc = map[string]string{}
	}
	group.RecipeVarsDoc[toks[1].content] = toks[2].content
	return instrs, p.nextLine(), nil
}
