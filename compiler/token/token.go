package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "//"

	IDENT = "IDENT"
	INT   = "INT"

	ADD = "+"
	SUB = "-"
	MUL = "*"
	QUO = "/"
	REM = "%"

	SHL = "<<"
	SHR = ">>"

	RANGE = ".."

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
	COLON     = ":"
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	VAR      = "var"
	RETURN   = "return"
	FUNCTION = "func"
	IF       = "if"
	ELSE     = "else"
	WHILE    = "while"
	FOR      = "for"
	IN       = "in"
)

var keywords = map[string]TokenType{
	"var":    VAR,
	"return": RETURN,
	"func":   FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"in":     IN,
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
