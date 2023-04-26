package lexer

import (
	"github.com/dsxg666/snakecoin/vm/token"
)

type Lexer struct {
	position     int    // current character position
	nextPosition int    // next character position
	character    rune   // current character
	characters   []rune // rune slice of input string
}

func New(input string) *Lexer {
	l := &Lexer{characters: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.characters) {
		l.character = rune(0)
	} else {
		l.character = l.characters[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition++
}

// NextToken to read next token, skipping the white space.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	// skip single-line comments
	// Unless they are immediately followed by a number, because
	// registers are "#N".
	if l.character == '#' {
		if !isDigit(l.peekChar()) {
			l.skipComment()
			return l.NextToken()
		}
	}

	switch l.character {
	case ',':
		tok = newToken(token.COMMA, l.character)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case ':':
		tok.Type = token.LABEL
		tok.Literal = l.readLabel()
	case rune(0):
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isDigit(l.character) {
			return l.readDecimal()
		}
		tok.Literal = l.readIdentifier()
		tok.Type = token.LookupIdentifier(tok.Literal)
		return tok
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.Type, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifier(l.character) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.character) {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.character != '\n' && l.character != rune(0) {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isHexDigit(l.character) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

func (l *Lexer) readUntilWhitespace() string {
	position := l.position
	for !isWhitespace(l.character) && l.character != rune(0) {
		l.readChar()
	}
	return string(l.characters[position:l.position])
}

func (l *Lexer) readDecimal() token.Token {
	integer := l.readNumber()
	if isEmpty(l.character) || isWhitespace(l.character) || l.character == ',' {
		return token.Token{Type: token.INT, Literal: integer}
	}

	illegalPart := l.readUntilWhitespace()
	return token.Token{Type: token.ILLEGAL, Literal: integer + illegalPart}
}

func (l *Lexer) readString() string {
	var out string
	for {
		l.readChar()
		if l.character == '"' {
			break
		}
		if l.character == '\\' {
			l.readChar()

			if l.character == 'n' {
				l.character = '\n'
			}
			if l.character == 'r' {
				l.character = '\r'
			}
			if l.character == 't' {
				l.character = '\t'
			}
			if l.character == '"' {
				l.character = '"'
			}
			if l.character == '\\' {
				l.character = '\\'
			}
		}
		out = out + string(l.character)
	}
	return out
}

func (l *Lexer) readLabel() string {
	return l.readUntilWhitespace()
}

func (l *Lexer) peekChar() rune {
	if l.nextPosition >= len(l.characters) {
		return rune(0)
	}
	return l.characters[l.nextPosition]
}

func isIdentifier(ch rune) bool {
	return ch != ',' && !isWhitespace(ch) && !isEmpty(ch)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isEmpty(ch rune) bool {
	return rune(0) == ch
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch rune) bool {
	if isDigit(ch) {
		return true
	}
	if 'a' <= ch && ch <= 'f' {
		return true
	}
	if 'A' <= ch && ch <= 'F' {
		return true
	}
	if ('x' == ch) || ('X' == ch) {
		return true
	}
	return false
}
