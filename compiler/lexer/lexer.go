package lexer

import (
	"fmt"
	"github.com/arata-nvm/Solitude/compiler/errors"
	"github.com/arata-nvm/Solitude/compiler/token"
	"io/ioutil"
)

type Lexer struct {
	filename     string
	input        string
	position     int
	readPosition int
	line         int
	ch           byte
}

func NewFromString(input string) *Lexer {
	l := &Lexer{
		filename: "__input__",
		input:    input,
		line:     1,
	}

	l.readChar()
	return l
}

func NewFromFile(filename string) (*Lexer, error) {
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	l := &Lexer{
		filename: filename,
		input:    string(code),
		line:     1,
	}

	l.readChar()
	return l, nil
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.New(tokenType, literal, l.getCurrentPos())
}

func (l *Lexer) getCurrentPos() token.Position {
	return token.Position{
		Filename: l.filename,
		Line:     l.line,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()

	var tok token.Token

	switch l.ch {
	case 0:
		tok = l.newToken(token.EOF, "")
	case '+':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.ADD_ASSIGN, "+=")
		} else {
			tok = l.newToken(token.ADD, "+")
		}
	case '-':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.SUB_ASSIGN, "-=")
		} else {
			tok = l.newToken(token.SUB, "-")
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.MUL_ASSIGN, "*=")
		} else {
			tok = l.newToken(token.MUL, "*")
		}
	case '/':
		if l.peekChar() == '/' {
			comment := l.readLine()
			tok = l.newToken(token.COMMENT, comment)
			break
		}

		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.QUO_ASSIGN, "/=")
		} else {
			tok = l.newToken(token.QUO, "/")
		}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.REM_ASSIGN, "%=")
		} else {
			tok = l.newToken(token.REM, "%")
		}
	case ',':
		tok = l.newToken(token.COMMA, ",")
	case ':':
		tok = l.newToken(token.COLON, ":")
	case ';':
		tok = l.newToken(token.SEMICOLON, ";")
	case '(':
		tok = l.newToken(token.LPAREN, "(")
	case ')':
		tok = l.newToken(token.RPAREN, ")")
	case '[':
		tok = l.newToken(token.LBRACKET, "[")
	case ']':
		tok = l.newToken(token.RBRACKET, "]")
	case '{':
		tok = l.newToken(token.LBRACE, "{")
	case '}':
		tok = l.newToken(token.RBRACE, "}")
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.EQ, "==")
		} else {
			tok = l.newToken(token.ASSIGN, "=")
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = l.newToken(token.NEQ, "!=")
		}
	case '<':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = l.newToken(token.LTE, "<=")
		case '<':
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.newToken(token.SHL_ASSIGN, "<<=")
			} else {
				tok = l.newToken(token.SHL, "<<")
			}
		default:
			tok = l.newToken(token.LT, "<")
		}
	case '>':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = l.newToken(token.GTE, ">=")
		case '>':
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.newToken(token.SHR_ASSIGN, ">>=")
			} else {
				tok = l.newToken(token.SHR, ">>")
			}
		default:
			tok = l.newToken(token.GT, ">")
		}
	case '.':
		if l.peekChar() == '.' {
			l.readChar()
			tok = l.newToken(token.RANGE, "..")
		} else {
			tok = l.newToken(token.PERIOD, ".")
		}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			t := token.LookUpIdent(ident)
			return l.newToken(t, ident)
		} else if isDigit(l.ch) {
			return l.readNumberLiteral()
		}
		errors.ErrorExit(fmt.Sprintf("%s | Illegal charactor '%c'", l.getCurrentPos(), l.ch))
	}

	l.readChar()

	return tok
}

func (l *Lexer) readNumberLiteral() token.Token {
	numLit := l.readNumber()
	// .. -> range
	if l.ch == '.' && l.peekChar() != '.' {
		// Float
		l.readChar()
		return l.newToken(token.FLOAT, fmt.Sprintf("%s.%s", numLit, l.readNumber()))
	}

	return l.newToken(token.INT, numLit)
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == '\n' || l.ch == '\r' {
		l.line++
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
	for isLetter(l.ch) || isDigit(l.ch) {
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
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
