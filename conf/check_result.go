package conf

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/tiramiseb/quickonf/commands"
)

type (
	CheckType     string
	CheckSeverity int
)

const (
	CheckTypeKeyword   CheckType = "keyword"
	CheckTypeGroupname CheckType = "groupname"
	CheckTypeVariable  CheckType = "variable"
	CheckTypeFunction  CheckType = "function"
	CheckTypeOperator  CheckType = "operator"
	CheckTypeUnknown   CheckType = "unknown"

	CheckSeverityError       = 0
	CheckSeverityWarning     = 1
	CheckSeverityInformation = 2
	CheckSeverityHint        = 3
)

type CheckToken struct {
	Content    string                `json:"content"`
	Help       string                `json:"help"`
	Completion []checkCompletionItem `json:"completion"`
	Line       int                   `json:"line"`
	Start      int                   `json:"start"`
	End        int                   `json:"end"`
	Length     int                   `json:"length"`
	Type       CheckType             `json:"type"`
}

type CheckError struct {
	Content  string        `json:"content"`
	Line     int           `json:"line"`
	Start    int           `json:"start"`
	End      int           `json:"end"`
	Length   int           `json:"length"`
	Severity CheckSeverity `json:"severity"`
	Message  string        `json:"message"`
}

type completionKind string

const (
	completionKindCommand  completionKind = "command"
	completionKindVariable completionKind = "variable"
)

type checkCompletionItem struct {
	Label checkCompletionItemLabel `json:"label"`
	Kind  completionKind           `json:"kind"`
}

type checkCompletionItemLabel struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type CheckResult struct {
	Tokens []CheckToken `json:"tokens"`
	Errors []CheckError `json:"errors"`
}

func newCheckResult() *CheckResult {
	return &CheckResult{
		Tokens: []CheckToken{},
		Errors: []CheckError{},
	}
}

func (r *CheckResult) addToken(tok *token, typ CheckType) {
	r.Tokens = append(r.Tokens, CheckToken{
		Content: tok.raw,
		Line:    tok.line - 1,
		Start:   tok.column - 1,
		End:     tok.column + tok.rawLength - 1,
		Length:  tok.rawLength,
		Type:    typ,
	})
}

func (r *CheckResult) addCommandToken(tok *token, cmd *commands.Command) {
	cmds := commands.ListStartWith(tok.content)
	completion := make([]checkCompletionItem, len(cmds))
	for i, cmd := range cmds {
		completion[i] = checkCompletionItem{
			Label: checkCompletionItemLabel{Label: cmd.Name},
			Kind:  completionKindCommand,
		}
	}
	r.Tokens = append(r.Tokens, CheckToken{
		Content:    tok.raw,
		Help:       cmd.MarkdownHelp(),
		Completion: completion,
		Line:       tok.line - 1,
		Start:      tok.column - 1,
		End:        tok.column + tok.rawLength - 1,
		Length:     tok.rawLength,
		Type:       CheckTypeFunction,
	})
}

func (r *CheckResult) addUnknownCommandToken(tok *token) {
	cmds := commands.ListStartWith(tok.content)
	completion := make([]checkCompletionItem, len(cmds))
	for i, cmd := range cmds {
		completion[i] = checkCompletionItem{
			Label: checkCompletionItemLabel{Label: cmd.Name},
			Kind:  completionKindCommand,
		}
	}
	r.Tokens = append(r.Tokens, CheckToken{
		Content:    tok.raw,
		Completion: completion,
		Line:       tok.line - 1,
		Start:      tok.column - 1,
		End:        tok.column + tok.rawLength - 1,
		Length:     tok.rawLength,
		Type:       CheckTypeUnknown,
	})
}

// makeVarsTokens adds tokens for vars inside a given token
func (r *CheckResult) makeVarsTokens(tok *token, knownVars map[string]string) {
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
				r.addUnfinishedVariableToken(&token{
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
				r.addVariableToken(&token{
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
		r.addUnfinishedVariableToken(&token{
			line:      tok.line + lineOffset,
			column:    tok.column + columnOffset,
			rawLength: currentVarLen + 1,
			raw:       currentVar,
		}, knownVars)
	}
}

func (r *CheckResult) addVariableToken(tok *token, knownVars map[string]string) {
	help := ""
	var found bool
	for key, instruction := range knownVars {
		if tok.raw == key {
			help = instruction
			found = true
		}
	}
	var completion []checkCompletionItem
	for key, instruction := range knownVars {
		if strings.HasPrefix(key, tok.content) {
			completion = append(completion, checkCompletionItem{
				Label: checkCompletionItemLabel{
					Label:       fmt.Sprintf("<%s>", key),
					Description: instruction,
				},
				Kind: completionKindVariable,
			})
		}
	}
	sort.Slice(completion, func(i, j int) bool {
		return completion[i].Label.Label < completion[j].Label.Label
	})
	r.Tokens = append(r.Tokens, CheckToken{
		Content:    tok.raw,
		Help:       help,
		Completion: completion,
		Line:       tok.line - 1,
		Start:      tok.column - 1,
		End:        tok.column + tok.rawLength - 1,
		Length:     tok.rawLength,
		Type:       CheckTypeVariable,
	})
	if !found {
		r.addError(tok, CheckSeverityInformation, "variable undefined, will not be translated")
	}
}

func (r *CheckResult) addUnfinishedVariableToken(tok *token, knownVars map[string]string) {
	var completion []checkCompletionItem
	for key, instruction := range knownVars {
		if strings.HasPrefix(key, tok.content) {
			completion = append(completion, checkCompletionItem{
				Label: checkCompletionItemLabel{
					Label:       fmt.Sprintf("<%s>", key),
					Description: instruction,
				},
				Kind: completionKindVariable,
			})
		}
	}
	sort.Slice(completion, func(i, j int) bool {
		return completion[i].Label.Label < completion[j].Label.Label
	})
	r.Tokens = append(r.Tokens, CheckToken{
		Content:    tok.raw,
		Completion: completion,
		Line:       tok.line - 1,
		Start:      tok.column - 1,
		End:        tok.column + tok.rawLength - 1,
		Length:     tok.rawLength,
		Type:       CheckTypeUnknown,
	})
}

func (r *CheckResult) addError(tok *token, severity CheckSeverity, message string) {
	r.Errors = append(r.Errors, CheckError{
		Content:  tok.raw,
		Line:     tok.line - 1,
		Start:    tok.column - 1,
		End:      tok.column + tok.rawLength - 1,
		Length:   tok.rawLength,
		Severity: severity,
		Message:  message,
	})
}

func (r *CheckResult) addErrorf(tok *token, severity CheckSeverity, format string, args ...interface{}) {
	r.addError(tok, severity, fmt.Sprintf(format, args...))
}

func (r *CheckResult) sort() {
	r.sortAndUniqTokens()
	r.sortAndUniqErrors()
}

func (r *CheckResult) sortAndUniqTokens() {
	sort.Slice(r.Tokens, func(i int, j int) bool {
		if r.Tokens[i].Line != r.Tokens[j].Line {
			return r.Tokens[i].Line < r.Tokens[j].Line
		}
		if r.Tokens[i].Start == r.Tokens[j].Start {
			return r.Tokens[i].Length <= r.Tokens[j].Length
		}
		return r.Tokens[i].Start < r.Tokens[j].Start
	})

	var (
		uniqTokens []CheckToken
		prev       CheckToken
	)
	for _, tok := range r.Tokens {
		if tok.Line != prev.Line ||
			tok.Start != prev.Start ||
			tok.End != prev.End ||
			tok.Length != prev.Length {
			uniqTokens = append(uniqTokens, tok)
			prev = tok
		}
	}
	r.Tokens = uniqTokens
}

func (r *CheckResult) sortAndUniqErrors() {
	sort.Slice(r.Errors, func(i int, j int) bool {
		if r.Errors[i].Line != r.Errors[j].Line {
			return r.Errors[i].Line < r.Errors[j].Line
		}
		if r.Errors[i].Start == r.Errors[j].Start {
			return r.Errors[i].Length <= r.Errors[j].Length
		}
		return r.Errors[i].Start < r.Errors[j].Start
	})

	var (
		uniqErrors []CheckError
		prev       CheckError
	)
	for _, err := range r.Errors {
		if err.Line != prev.Line ||
			err.Start != prev.Start ||
			err.End != prev.End ||
			err.Length != prev.Length {
			uniqErrors = append(uniqErrors, err)
			prev = err
		}
	}
	r.Errors = uniqErrors
}

func (r *CheckResult) ToJSON() (string, error) {
	result, err := json.MarshalIndent(r, "", "  ")
	return string(result), err
}
