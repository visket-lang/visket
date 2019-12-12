package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	INT = "INT"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	EQ  = "=="
	NEQ = "!="

	LPAREN = "("
	RPAREN = ")"
)

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal string) Token {
	tok := Token{
		Type:    tokenType,
		Literal: literal,
	}

	return tok
}
