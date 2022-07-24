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

	curLine   int
	curIndent int
	curCol    int

	currentWord   []byte
	currentRaw    []byte
	curWordLine   int
	curWordCol    int
	curWordLength int
	curRawLength  int
}

type lexerContext int

const (
	contextStartOfLine lexerContext = iota // About to read the first character of a line
	contextIndentation                     // Reading spaces from the beginning of the line
	contextComment                         // Reading comment (starting with "#")
	contextGroupName                       // Reading the name of a group
	contextDefault                         // Reading some word, nothing special
	contextQuotes                          // Reading the content of a quoted string
	contextCookbookURI                     // Reading the URI of a recipes book
	contextSpace                           // Reading a space
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
		case contextCookbookURI:
			currentContext, err = l.cookbook()
		case contextSpace:
			currentContext, err = l.space()
		}
	}
	if errors.Is(err, io.EOF) {
		l.tokens = append(l.tokens, l.eolToken())
		err = nil
	}
	return l.tokens, err
}

func (l *lexer) next() (byte, error) {
	l.curIndent++
	l.curCol++
	return l.r.ReadByte()
}

func (l *lexer) resetWord() {
	l.curWordLine = l.curLine
	l.curWordCol = l.curCol
	l.curWordLength = 0
	l.curRawLength = 0
	l.currentWord = l.currentWord[:0]
	l.currentRaw = l.currentRaw[:0]
}

func (l *lexer) resetWordWith(b byte) {
	l.curWordLine = l.curLine
	l.curWordCol = l.curCol
	l.curWordLength = 1
	l.curRawLength = 1
	l.currentWord = l.currentWord[:0]
	l.currentRaw = l.currentRaw[:0]
	l.currentWord = append(l.currentWord, b)
	l.currentRaw = append(l.currentRaw, b)
}

func (l *lexer) resetWordWithQuote() {
	l.curWordLine = l.curLine
	l.curWordCol = l.curCol
	l.curWordLength = 1
	l.curRawLength = 1
	l.currentWord = l.currentWord[:0]
	l.currentRaw = l.currentRaw[:0]
	l.currentRaw = append(l.currentRaw, '"')
}

func (l *lexer) resetWordWithEscape() {
	l.curWordLine = l.curLine
	l.curWordCol = l.curCol
	l.curWordLength = 1
	l.curRawLength = 1
	l.currentWord = l.currentWord[:0]
	l.currentRaw = l.currentRaw[:0]
	l.currentRaw = append(l.currentRaw, '\\')
}

func (l *lexer) appendToWord(b byte) {
	l.curWordLength++
	l.curRawLength++
	l.currentWord = append(l.currentWord, b)
	l.currentRaw = append(l.currentRaw, b)
}

func (l *lexer) appendEscapeToRaw() {
	l.curRawLength++
	l.currentRaw = append(l.currentRaw, '\\')
}

func (l *lexer) eolToken() *token {
	return &token{l.curLine, l.curCol, 0, 0, tokenEOL, "", ""}
}

func (l *lexer) identifyToken() *token {
	return identifyToken(l.curWordLine, l.curWordCol, l.curWordLength, l.curRawLength, string(l.currentWord), string(l.currentRaw))
}

func (l *lexer) indentToken() *token {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(l.curIndent-1))
	return &token{l.curLine, 1, l.curCol - 1, l.curCol - 1, tokenIndentation, string(buf), string(buf)}
}

func (l *lexer) defaultToken() *token {
	return &token{l.curWordLine, l.curWordCol, l.curWordLength, l.curRawLength, tokenDefault, strings.TrimSpace(string(l.currentWord)), string(l.currentRaw)}
}

func (l *lexer) startOfLine() (lexerContext, error) {
	l.curLine++
	l.curIndent = 0
	l.curCol = 0
	b, err := l.next()
	if err != nil {
		return contextStartOfLine, err
	}
	switch b {
	case ' ':
	case '\t':
		l.curIndent = 8
	case '\n':
		return contextStartOfLine, nil
	case '\\':
		l.resetWordWithEscape()
		b, err := l.next()
		if err != nil {
			return contextStartOfLine, err
		}
		l.appendToWord(b)
		return contextGroupName, nil
	case '#':
		return contextComment, nil
	default:
		l.resetWordWith(b)
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
			l.curIndent = l.curIndent + 7 - (l.curIndent % 8)
		case '\n':
			return contextStartOfLine, nil
		case '#':
			return contextComment, nil
		case '"':
			l.tokens = append(l.tokens, l.indentToken())
			l.resetWordWithQuote()
			return contextQuotes, nil
		case '\\':
			l.tokens = append(l.tokens, l.indentToken())
			l.resetWordWithEscape()
			b, err := l.next()
			if err != nil {
				return contextStartOfLine, err
			}
			l.appendToWord(b)
			return contextDefault, nil
		default:
			l.tokens = append(l.tokens, l.indentToken())
			l.resetWordWith(b)
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
	for {
		b, err := l.next()
		if err != nil {
			return contextGroupName, err
		}
		switch b {
		case ' ', '\t':
			// At the first space, check if it is something else than a group name
			switch string(l.currentWord) {
			case "cookbook":
				// Prepare to get the cookbook URI
				l.tokens = append(l.tokens,
					&token{l.curWordLine, l.curWordCol, l.curWordLength, l.curRawLength, tokenCookbook, "cookbook", "cookbook"},
				)
				l.resetWord()
				l.curWordCol++
				return contextCookbookURI, nil
			}
			l.appendToWord(b)
		case '\n':
			l.tokens = append(l.tokens,
				&token{l.curLine, 1, l.curWordLength, l.curRawLength, tokenGroupName, strings.TrimSpace(string(l.currentWord)), string(l.currentRaw)},
				l.eolToken(),
			)
			return contextStartOfLine, nil
		case '#':
			l.tokens = append(l.tokens,
				&token{l.curLine, 1, l.curWordLength, l.curRawLength, tokenGroupName, strings.TrimSpace(string(l.currentWord)), string(l.currentRaw)},
				l.eolToken(),
			)
			return contextComment, nil
		case '\\':
			l.appendEscapeToRaw()
			b, err := l.next()
			if err != nil {
				return contextGroupName, err
			}
			l.appendToWord(b)
		default:
			l.appendToWord(b)
		}
	}
}

func (l *lexer) defaut() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				l.tokens = append(l.tokens, l.identifyToken())
			}
			return contextDefault, err
		}
		switch b {
		case '\n':
			l.tokens = append(l.tokens, l.identifyToken(), l.eolToken())
			return contextStartOfLine, nil
		case ' ', '\t':
			l.tokens = append(l.tokens, l.identifyToken())
			return contextSpace, nil
		case '#':
			l.tokens = append(l.tokens, l.identifyToken(), l.eolToken())
			return contextComment, nil
		case '"':
			return contextQuotes, nil
		case '\\':
			l.appendEscapeToRaw()
			b, err := l.next()
			if err != nil {
				return contextDefault, err
			}
			l.appendToWord(b)
		default:
			l.appendToWord(b)
		}
	}
}

func (l *lexer) quotes() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = errors.New("unclosed quote at end of file")
			}
			return contextQuotes, err
		}
		switch b {
		case '"':
			return contextDefault, nil
		case '\n':
			l.curLine++
			l.curCol = 0
			l.appendToWord(b)
		case '\\':
			l.appendEscapeToRaw()
			b, err := l.next()
			if err != nil {
				return contextQuotes, err
			}
			if b == '\n' {
				l.curLine++
				l.curCol = 0
			}
			l.appendToWord(b)
		default:
			l.appendToWord(b)
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
			// Nothing to do, it is a space
		case '\n':
			l.tokens = append(l.tokens, l.eolToken())
			return contextStartOfLine, nil
		case '#':
			l.tokens = append(l.tokens, l.eolToken())
			return contextComment, nil
		case '"':
			l.resetWordWithQuote()
			return contextQuotes, nil
		case '\\':
			l.resetWordWithEscape()
			b, err := l.next()
			if err != nil {
				return contextSpace, err
			}
			l.appendToWord(b)
			return contextDefault, nil
		default:
			l.resetWordWith(b)
			return contextDefault, nil
		}
	}
}

func (l *lexer) cookbook() (lexerContext, error) {
	for {
		b, err := l.next()
		if err != nil {
			return contextCookbookURI, err
		}
		switch b {
		case '\n':
			l.tokens = append(l.tokens, l.defaultToken(), l.eolToken())
			return contextStartOfLine, nil
		case '#':
			l.tokens = append(l.tokens, l.defaultToken(), l.eolToken())
			return contextComment, nil
		case '\\':
			b, err := l.next()
			if err != nil {
				return contextCookbookURI, err
			}
			l.appendToWord(b)
		default:
			l.appendToWord(b)
		}
	}
}
