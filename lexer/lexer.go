package lexer

import (
	"github.com/arata-nvm/Solitude/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()

	var tok token.Token

	switch l.ch {
	case 0:
		tok = token.New(token.EOF, "")
	case '+':
		tok = token.New(token.PLUS, "+")
	case '-':
		tok = token.New(token.MINUS, "-")
	case '*':
		tok = token.New(token.ASTERISK, "*")
	case '/':
		tok = token.New(token.SLASH, "/")
	case '(':
		tok = token.New(token.LPAREN, "(")
	case ')':
		tok = token.New(token.RPAREN, ")")
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.EQ, "==")
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.NEQ, "!=")
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LTE, "<=")
		} else {
			tok = token.New(token.LT, "<")
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GTE, ">=")
		} else {
			tok = token.New(token.GT, ">")
		}
	default:
		if isDigit(l.ch) {
			numLit := l.readNumber()
			return token.New(token.INT, numLit)
		} else {
			return token.New(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readNumber() string {
	readPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[readPos:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}