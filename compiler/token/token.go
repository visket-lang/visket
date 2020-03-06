package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "//"

	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	ADD = "+"
	SUB = "-"
	MUL = "*"
	QUO = "/"
	REM = "%"

	SHL = "<<"
	SHR = ">>"

	RANGE  = ".."
	MODSEP = "::"

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
	PERIOD    = "."
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
	STRUCT   = "struct"
	NEW      = "new"
	IMPORT   = "import"
	REF      = "ref"
	VAL      = "val"
	MODULE   = "module"
	INCLUDE  = "include"
)

var keywords = map[string]TokenType{
	"var":     VAR,
	"return":  RETURN,
	"func":    FUNCTION,
	"if":      IF,
	"else":    ELSE,
	"while":   WHILE,
	"for":     FOR,
	"in":      IN,
	"struct":  STRUCT,
	"new":     NEW,
	"import":  IMPORT,
	"ref":     REF,
	"val":     VAL,
	"module":  MODULE,
	"include": INCLUDE,
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

func New(tokenType TokenType, literal string, pos Position) Token {
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
