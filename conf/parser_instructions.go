package conf

import "github.com/tiramiseb/quickonf/instructions"

func (p *parser) instructions(prefixAllWith, line tokens, group *instructions.Group, currentIndent int) (instrs []instructions.Instruction, next tokens) {
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
			p.errs = append(p.errs, toks[0].errorf("invalid indentation (expecting %d)", currentIndent))
			return instrs, p.nextLine()
		case thisIndent < currentIndent:
			// This indentation block is finished, quit the function
			return instrs, line
		}

		// Add prefix to instruction if needed
		toks = addPrefix(prefixAllWith, toks)

		// Parse the tokens for this instruction!
		var ins []instructions.Instruction
		switch toks[0].typ {
		case tokenExpand:
			ins, line = p.expand(toks)
		case tokenIf:
			ins, line = p.ifThen(toks, group, currentIndent)
		case tokenPriority:
			ins, line = p.priority(toks, group)
		case tokenRecipe:
			ins, line = p.recipe(toks, group, currentIndent)
		case tokenDoc:
			ins, line = p.recipeDoc(toks, group)
		case tokenVardoc:
			ins, line = p.recipeVarDoc(toks, group)
		case tokenRepeat:
			ins, line = p.repeat(toks, group, currentIndent)
		default:
			ins, line = p.command(toks)
		}
		instrs = append(instrs, ins...)
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
