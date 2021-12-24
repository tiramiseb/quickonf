package conf

import (
	"encoding/binary"
	"fmt"
)

type tokenType int

const (
	tokenGroupName tokenType = iota
	tokenIndentation
	tokenEOL
	tokenDefault

	tokenIf
	tokenEqual
)

var tokenOperators = []tokenType{tokenEqual}

type token struct {
	// Position of the first character of the token (for debugging)
	line   int
	column int

	typ     tokenType
	content string
}

type tokens []*token

func identifyToken(line int, column int, content string) *token {
	typ := tokenDefault
	switch content {
	case "if":
		typ = tokenIf
	case "=":
		typ = tokenEqual
	}
	return &token{line, column, typ, content}
}

func (t *token) isOperator() bool {
	for _, ope := range tokenOperators {
		if t.typ == ope {
			return true
		}
	}
	return false
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
	size, _ := binary.Uvarint([]byte(t[0].content))
	return int(size), t[1]
}