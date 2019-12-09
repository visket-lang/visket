package token

import (
	"log"
	"strconv"
)

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	INT               = "INT"
	PLUS              = "PLUS"
	MINUS             = "MINUS"
	EOF               = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
	Val     int
}

func New(tokenType TokenType, literal string) Token {
	tok := Token{
		Type:    tokenType,
		Literal: literal,
	}

	if tok.Type == INT {
		n, err := strconv.Atoi(tok.Literal)
		if err != nil {
			log.Fatal(err)
		}
		tok.Val = n
	}

	return tok
}
