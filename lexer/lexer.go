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
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.ADD_ASSIGN, "+=")
		} else {
			tok = token.New(token.ADD, "+")
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.SUB_ASSIGN, "-=")
		} else {
			tok = token.New(token.SUB, "-")
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.MUL_ASSIGN, "*=")
		} else {
			tok = token.New(token.MUL, "*")
		}
	case '/':
		if l.peekChar() == '/' {
			comment := l.readLine()
			tok = token.New(token.COMMENT, comment)
			break
		}

		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.QUO_ASSIGN, "/=")
		} else {
			tok = token.New(token.QUO, "/")
		}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.REM_ASSIGN, "%=")
		} else {
			tok = token.New(token.REM, "%")
		}
	case ',':
		tok = token.New(token.COMMA, ",")
	case ';':
		tok = token.New(token.SEMICOLON, ";")
	case '(':
		tok = token.New(token.LPAREN, "(")
	case ')':
		tok = token.New(token.RPAREN, ")")
	case '{':
		tok = token.New(token.LBRACE, "{")
	case '}':
		tok = token.New(token.RBRACE, "}")
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.EQ, "==")
		} else {
			tok = token.New(token.ASSIGN, "=")
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.NEQ, "!=")
		}
	case '<':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.New(token.LTE, "<=")
		case '<':
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.New(token.SHL_ASSIGN, "<<=")
			} else {
				tok = token.New(token.SHL, "<<")
			}
		default:
			tok = token.New(token.LT, "<")
		}
	case '>':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = token.New(token.GTE, ">=")
		case '>':
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = token.New(token.SHR_ASSIGN, ">>=")
			} else {
				tok = token.New(token.SHR, ">>")
			}
		default:
			tok = token.New(token.GT, ">")
		}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			t := token.LookUpIdent(ident)
			return token.New(t, ident)
		} else if isDigit(l.ch) {
			numLit := l.readNumber()
			return token.New(token.INT, numLit)
		}
		return token.New(token.ILLEGAL, string(l.ch))
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

func (l *Lexer) readIdentifier() string {
	readPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[readPos:l.position]
}

func (l *Lexer) readNumber() string {
	readPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[readPos:l.position]
}

func (l *Lexer) readLine() string {
	readPos := l.position
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
			break
		}
	}

	return l.input[readPos:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
