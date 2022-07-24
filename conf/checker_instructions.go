package conf

func (c *checker) instructions(prefixAllWith, line tokens, currentIndent int, knownVars []string) (next tokens, newVars []string) {
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
			c.result.addErrorf(toks[0], CheckSeverityError, "invalid indentation (expecting %d)", currentIndent)
			return c.nextLine(), newVars
		case thisIndent < currentIndent:
			// This indentation block is finished, quit the function
			return line, newVars
		}

		// Add prefix to instruction if needed
		toks = addPrefix(prefixAllWith, toks)

		var thisNewVars []string
		// Parse the tokens for this instruction!
		switch toks[0].typ {
		case tokenExpand:
			line, thisNewVars = c.expand(toks, knownVars)
		case tokenIf:
			line, thisNewVars = c.ifThen(toks, currentIndent, knownVars)
		case tokenPriority:
			line, thisNewVars = c.priority(toks)
		case tokenRecipe:
			line, thisNewVars = c.recipe(toks, currentIndent, knownVars)
		case tokenDoc:
			line, thisNewVars = c.recipeDoc(toks)
		case tokenVardoc:
			line, thisNewVars = c.recipeVarDoc(toks)
		case tokenRepeat:
			line, thisNewVars = c.repeat(toks, currentIndent, knownVars)
		default:
			line, thisNewVars = c.command(toks, knownVars)
		}
		newVars = append(newVars, thisNewVars...)
		knownVars = append(knownVars, thisNewVars...)
	}
}
