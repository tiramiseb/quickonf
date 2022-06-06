package conf

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

type lexer struct {
	r *bufio.Reader

	tokens tokens

	curLine int
	curCol  int

	currentWord []byte
	curWordLine int
	curWordCol  int
}

type lexerContext int

const (
	contextStartOfLine lexerContext = iota
	contextIndentation
	contextComment
	contextGroupName
	contextDefault
	contextQuotes
	contextSpace
)

func newLexer(r io.Reader) *lexer {
	return &lexer{
		r: bufio.NewReader(r),
	}
}

// scan is the lexer, it transforms an io.Reader to a list of tokens
func (l *lexer) scan() (tokens, error) {
	var err error
	currentContext := contextStartOfLine
	for err == nil {
		switch currentContext {
		case contextStartOfLine:
			currentContext, err = l.startOfLine()
		case contextIndentation:
			currentContext, err = l.indentation()
		case contextComment:
			currentContext, err = l.comment()
		case contextGroupName:
			currentContext, err = l.groupName()
		case contextDefault:
			currentContext, err = l.defaut()
		case contextQuotes:
			currentContext, err = l.quotes()
		case contextSpace:
			currentContext, err = l.space()
		}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return l.tokens, err
}

func (l *lexer) next() (byte, error) {
	l.curCol++
	return l.r.ReadByte()
}

func (l *lexer) startOfLine() (lexerContext, error) {
	l.curLine++
	l.curCol = 0
	b, err := l.next()
	if err != nil {
		return contextStartOfLine, err
	}
	switch b {
	case ' ':
	case '\t':
		l.curCol = 8
	case '\n':
		return contextStartOfLine, nil
	case '\\':
		l.curWordLine = l.curLine
		l.curWordCol = l.curCol
		return contextGroupName, nil
	case '#':
		return contextComment, nil
	default:
		l.curWordLine = l.curLine
		l.curWordCol = l.curCol
		l.currentWord = []byte{b}
		return contextGroupName, nil
	}
	return contextIndentation, nil
}

func (l *lexer) indentation() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextStartOfLine, err
		}
		switch b {
		case ' ':
		case '\t':
			l.curCol = l.curCol + 8 - (l.curCol % 8)
		case '\n':
			return contextStartOfLine, nil
		case '#':
			return contextComment, nil
		case '"':
			buf := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(buf, uint64(l.curCol-1))
			l.tokens = append(l.tokens, &token{l.curLine, 1, tokenIndentation, string(buf)})
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			l.currentWord = l.currentWord[:0]
			return contextQuotes, nil
		case '\\':
			buf := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(buf, uint64(l.curCol-1))
			l.tokens = append(l.tokens, &token{l.curLine, 1, tokenIndentation, string(buf)})
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			return contextDefault, nil
		default:
			buf := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(buf, uint64(l.curCol-1))
			l.tokens = append(l.tokens, &token{l.curLine, 1, tokenIndentation, string(buf)})
			l.currentWord = []byte{b}
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			return contextDefault, nil
		}
	}
}

func (l *lexer) comment() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextComment, err
		}
		if b == '\n' {
			return contextStartOfLine, nil
		}
	}
}

func (l *lexer) groupName() (lexerContext, error) {
	b, err := l.next()
	if err != nil {
		return contextGroupName, err
	}
	l.currentWord = append(l.currentWord, b)
	for {
		b, err := l.next()
		if err != nil {
			return contextGroupName, err
		}
		switch b {
		case '\n':
			l.tokens = append(l.tokens, &token{
				l.curLine, 1, tokenGroupName,
				strings.TrimSpace(string(l.currentWord)),
			})
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextStartOfLine, nil
		case '#':
			l.tokens = append(l.tokens, &token{
				l.curLine, 1, tokenGroupName,
				strings.TrimSpace(string(l.currentWord)),
			})
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextComment, nil
		case '\\':
			b, err := l.next()
			if err != nil {
				return contextGroupName, err
			}
			l.currentWord = append(l.currentWord, b)
		default:
			l.currentWord = append(l.currentWord, b)
		}
	}
}

func (l *lexer) defaut() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextDefault, err
		}
		switch b {
		case '\n':
			l.tokens = append(l.tokens, identifyToken(l.curWordLine, l.curWordCol, string(l.currentWord)))
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextStartOfLine, nil
		case ' ', '\t':
			l.tokens = append(l.tokens, identifyToken(l.curWordLine, l.curWordCol, string(l.currentWord)))
			return contextSpace, nil
		case '#':
			l.tokens = append(l.tokens, identifyToken(l.curWordLine, l.curWordCol, string(l.currentWord)))
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextComment, nil
		case '"':
			return contextQuotes, nil
		case '\\':
			b, err := l.next()
			if err != nil {
				return contextDefault, err
			}
			l.currentWord = append(l.currentWord, b)
		default:
			l.currentWord = append(l.currentWord, b)
		}
	}
}

func (l *lexer) quotes() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextQuotes, err
		}
		switch b {
		case '"':
			return contextDefault, nil
		case '\n':
			l.curLine++
			l.curCol = 0
			l.currentWord = append(l.currentWord, b)
		case '\\':
			b, err := l.next()
			if err != nil {
				return contextQuotes, err
			}
			if b == '\n' {
				l.curLine++
				l.curCol = 0
			}
			l.currentWord = append(l.currentWord, b)
		default:
			l.currentWord = append(l.currentWord, b)
		}
	}
}

func (l *lexer) space() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextSpace, err
		}
		switch b {
		case ' ', '\t':
		case '\n':
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextStartOfLine, nil
		case '#':
			l.tokens = append(l.tokens, &token{l.curLine, l.curCol, tokenEOL, ""})
			return contextComment, nil
		case '"':
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			l.currentWord = l.currentWord[:0]
			return contextQuotes, nil
		case '\\':
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			return contextDefault, nil
		default:
			l.currentWord = []byte{b}
			l.curWordLine = l.curLine
			l.curWordCol = l.curCol
			return contextDefault, nil
		}
	}
}
