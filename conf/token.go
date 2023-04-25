package conf

import (
	"encoding/binary"
	"fmt"

	"github.com/tiramiseb/quickonf/instructions"
)

type tokenType int

const (
	tokenGroupName tokenType = iota
	tokenIndentation
	tokenEOL
	tokenDefault

	tokenEqual
	tokenExpand
	tokenIf
	tokenPriority
	tokenCookbook
	tokenRecipe
	tokenDoc
	tokenVardoc
	tokenRepeat
)

type token struct {
	// Position of the first character of the token (for debugging)
	line      int
	column    int
	length    int
	rawLength int

	typ     tokenType
	content string
	raw     string
}

type tokens []*token

func identifyToken(line, column, length, rawLength int, content, raw string) *token {
	typ := tokenDefault
	switch content {
	case "=":
		typ = tokenEqual
	case "expand":
		typ = tokenExpand
	case "if":
		typ = tokenIf
	case "priority":
		typ = tokenPriority
	case "recipe":
		typ = tokenRecipe
	case "doc":
		typ = tokenDoc
	case "vardoc":
		typ = tokenVardoc
	case "repeat":
		typ = tokenRepeat
	default:
		typ = tokenDefault
	}
	return &token{line, column, length, rawLength, typ, content, raw}
}

func (t *token) error(msg string) *instructions.Error {
	return &instructions.Error{
		Line:    t.line,
		Column:  t.column,
		Message: msg,
	}
}

func (t *token) errorf(format string, a ...interface{}) *instructions.Error {
	return &instructions.Error{
		Line:    t.line,
		Column:  t.column,
		Message: fmt.Sprintf(format, a...),
	}
}

// indentations returns the indentation size and the remaining token of the line
func (t tokens) indentation() (int, tokens) {
	if len(t) == 0 {
		return 0, nil
	}
	if t[0].typ != tokenIndentation {
		return 0, t
	}
	size, _ := binary.Uvarint([]byte(t[0].content))
	return int(size), t[1:]
}
