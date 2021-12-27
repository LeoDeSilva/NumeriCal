package lexer

import "errors"

type Token struct {
	Type    string
	Literal string
}

func NewToken(tokenType string, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func lookupIdentifier(identifier string) (string, error) {
	if len(identifier) == 0 {
		return "", errors.New("lookupIdentifier: StringLength must be greater than 0")
	}

	if token, ok := keywords[identifier]; ok {
		return token, nil
	}

	return IDENTIFIER, nil
}

var keywords = map[string]string{
	"in":     IN,
	"define": DEFINE,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INT        = "INT"
	STRING     = "STRING"
	FLOAT      = "FLOAT"
	IDENTIFIER = "IDENTIFIER"
	UNIT       = "UNIT"

	EQ  = "EQ"
	EE  = "EE"
	NE  = "NE"
	LT  = "LT"
	LTE = "LTE"
	GT  = "GT"
	GTE = "GTE"

	ADD = "ADD"
	SUB = "SUB"
	DIV = "DIV"
	MUL = "MUL"
	POW = "POW"
	MOD = "MOD"

	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	LSQUARE = "LSQUARE"
	RSQUARE = "RSQUARE"
	LBRACE  = "LBRACE"
	RBRACE  = "RBRACE"

	SEMICOLON = "SEMICOLON"
	ARROW     = "ARROW"

	IN     = "IN"
	DEFINE = "DEFINE"
)
