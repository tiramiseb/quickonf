package conf

import (
	"github.com/tiramiseb/quickonf/instructions"
)

func (p *parser) recipeDoc(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens) {
	if len(toks) != 2 {
		p.errs = append(p.errs, toks[0].error("expected a documentation"))
	}
	group.RecipeDoc = toks[1].content
	return nil, p.nextLine()
}

func (p *parser) recipeVarDoc(toks tokens, group *instructions.Group) (instrs []instructions.Instruction, next tokens) {
	if len(toks) != 3 {
		p.errs = append(p.errs, toks[0].error("expected a variable name and a documentation"))
	}
	if group.RecipeVarsDoc == nil {
		group.RecipeVarsDoc = map[string]string{}
	}
	group.RecipeVarsDoc[toks[1].content] = toks[2].content
	return nil, p.nextLine()
}
