package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"
	COMMENT           = "//"

	IDENT = "IDENT"
	INT   = "INT"

	ADD = "+"
	SUB = "-"
	MUL = "*"
	QUO = "/"
	REM = "%"

	SHL = "<<"
	SHR = ">>"

	ADD_ASSIGN = "+="
	SUB_ASSIGN = "-="
	MUL_ASSIGN = "*="
	QUO_ASSIGN = "/="
	REM_ASSIGN = "%="

	SHL_ASSIGN = "<<="
	SHR_ASSIGN = ">>="

	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="

	ASSIGN = "="

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
	Pos     *Position
}

func New(tokenType TokenType, literal string, pos *Position) Token {
	tok := Token{
		Type:    tokenType,
		Literal: literal,
		Pos:     pos,
	}

	return tok
}

func LookUpIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}

	return IDENT
}
