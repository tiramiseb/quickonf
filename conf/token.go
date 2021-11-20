package conf

import "fmt"

type tokenType int

const (
	tokenGroupName tokenType = iota
	tokenIndentation
	tokenEOL
	tokenDefault
)

type token struct {
	// Position of the first character of the token (for debugging)
	line   int
	column int

	typ     tokenType
	content interface{}
}

type tokens []*token

func identifyToken(line int, column int, content interface{}) *token {
	// TODO Identify special tokens
	typ := tokenDefault
	content = string(content.([]uint8))
	return &token{line, column, typ, content}
}

func (t *token) error(msg string) error {
	return fmt.Errorf("[%d:%d] %s", t.line, t.column, msg)
}

func (t *token) errorf(format string, a ...interface{}) error {
	a = append([]interface{}{t.line, t.column}, a...)
	return fmt.Errorf("[%d:%d] "+format, a...)
}

// indentations returns the indentation size and the first token of the line
func (t tokens) indentation() (int, *token) {
	if len(t) == 0 {
		return 0, nil
	}
	if t[0].typ != tokenIndentation {
		return 0, t[0]
	}
	return t[0].content.(int), t[1]
}
