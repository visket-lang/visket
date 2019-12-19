package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	VAR      = "var"
	RETURN   = "return"
	FUNCTION = "func"
	IF       = "if"
	ELSE     = "else"
	WHILE    = "while"
	FOR      = "for"
)

var keywords = map[string]TokenType{
	"var":    VAR,
	"return": RETURN,
	"func":   FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
}

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

func LookUpIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}

	return IDENT
}
