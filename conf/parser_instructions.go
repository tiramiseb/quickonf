package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) instructions(prefixAllWith, line tokens, group *instructions.Group, currentIndent int, knownVars map[string]string) (instrs []instructions.Instruction, next tokens, newVars map[string]string) {
	newVars = map[string]string{}
	// Read a list of instructions...
	for {
		if len(line) == 0 {
			// End of file (the lexer doesn't leave any empty line)
			return
		}
		thisIndent, toks := line.indentation()

		switch {
		case thisIndent > currentIndent:
			// A larger indentation should not happen, unless we are in another block (which would be started when parsing an instruction needing another block)
			instrs = append(instrs, toks[0].errorf("invalid indentation (expecting %d)", currentIndent))
			p.checkResult.addErrorf(toks[0], CheckSeverityError, "invalid indentation (expecting %d)", currentIndent)
			return instrs, p.nextLine(), newVars
		case thisIndent < currentIndent:
			// This indentation block is finished, quit the function
			return instrs, line, newVars
		}

		// Add prefix to instruction if needed
		toks = addPrefix(prefixAllWith, toks)

		var thisNewVars map[string]string
		// Parse the tokens for this instruction!
		var ins []instructions.Instruction
		switch toks[0].typ {
		case tokenExpand:
			ins, line, thisNewVars = p.expand(toks, knownVars)
		case tokenIf:
			ins, line, thisNewVars = p.ifThen(toks, group, currentIndent, knownVars)
		case tokenPriority:
			ins, line, thisNewVars = p.priority(toks, group)
		case tokenRecipe:
			ins, line, thisNewVars = p.recipe(toks, group, currentIndent)
		case tokenDoc:
			ins, line, thisNewVars = p.recipeDoc(toks, group)
		case tokenVardoc:
			ins, line, thisNewVars = p.recipeVarDoc(toks, group)
		case tokenRepeat:
			ins, line, thisNewVars = p.repeat(toks, group, currentIndent, knownVars)
		default:
			ins, line, thisNewVars = p.command(toks, knownVars)
		}
		instrs = append(instrs, ins...)
		for k, v := range thisNewVars {
			newVars[k] = v
		}
		for k, v := range thisNewVars {
			knownVars[k] = v
		}
	}
}

func addPrefix(prefix tokens, existing tokens) tokens {
	if len(prefix) == 0 {
		return existing
	}

	nbPrefixes := len(prefix)

	newTokens := make(tokens, len(prefix)+len(existing))

	copy(newTokens, prefix)

	for i, t := range existing {
		newTokens[i+nbPrefixes] = t
	}

	return newTokens
}
