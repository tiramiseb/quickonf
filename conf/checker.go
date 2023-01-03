package conf

type checker struct {
	result *CheckResult

	tokens tokens
	idx    int
}

func newChecker(tokens tokens) checker {
	return checker{
		result: newCheckResult(),
		tokens: tokens,
	}
}

func (c *checker) nextLine() (toks tokens) {
	for c.idx < len(c.tokens) && c.tokens[c.idx].typ != tokenEOL {
		toks = append(toks, c.tokens[c.idx])
		c.idx++
	}
	// Ignore tokenEOL
	c.idx++
	return toks
}

// check parses the tokens in order to return a complete check result
//
// All functions called from this function receive the "next" line for
// processing and return the "next" line for processing by another sub-parser.
// It is necessary in order to know how to process next line
func (c *checker) check() *CheckResult {
	next := c.nextLine()
	for next != nil {
		next = c.noIndentation(next)
	}
	c.result.sort()
	return c.result
}

// makeVarsTokens adds tokens for vars inside a given token
func (c *checker) makeVarsTokens(tok *token, knownVars map[string]string) {
	var (
		currentVar    string
		currentVarLen int
		readingVar    bool
		columnOffset  int
		lineOffset    int
	)
	for i, char := range tok.raw {
		switch char {
		case '\n':
			lineOffset++
			columnOffset = 0 - tok.column
		case '<':
			if readingVar {
				c.addUnfinishedVarToken(&token{
					line:      tok.line + lineOffset,
					column:    tok.column + columnOffset,
					rawLength: currentVarLen + 1,
					raw:       currentVar,
				}, knownVars)
			}
			currentVar = ""
			columnOffset = i
			currentVarLen = 1
			readingVar = true
		case '>':
			if readingVar {
				c.addVarToken(&token{
					line:      tok.line + lineOffset,
					column:    tok.column + columnOffset,
					rawLength: currentVarLen + 1,
					raw:       currentVar,
				}, knownVars)
				readingVar = false
			}
		}
		if readingVar && char != '<' {
			currentVarLen++
			currentVar += string(char)
		}
	}
	if readingVar {
		c.addUnfinishedVarToken(&token{
			line:      tok.line + lineOffset,
			column:    tok.column + columnOffset,
			rawLength: currentVarLen + 1,
			raw:       currentVar,
		}, knownVars)
	}
}

func (c *checker) addUnfinishedVarToken(tok *token, knownVars map[string]string) {
	c.result.addUnfinishedVariableToken(tok, knownVars)
}

func (c *checker) addVarToken(tok *token, knownVars map[string]string) {
	c.result.addVariableToken(tok, knownVars)
	for key, instruction := range knownVars {
		if tok.raw == key {
			if instruction != "" {
				c.result.addError(tok, CheckSeverityInformation, instruction)
			}
			return
		}
	}
	c.result.addError(tok, CheckSeverityWarning, "variable undefined, will not be translated")
}
