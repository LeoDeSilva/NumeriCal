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
		return "", errors.New("LexerError: lookupIdentifier() identifier StringLength must be greater than 0")
	}

	if token, ok := keywords[identifier]; ok {
		return token, nil
	}

	return IDENTIFIER, nil
}

var keywords = map[string]string{
	"in":     IN,
	"define": DEFINE,
	"per":    DIV,
	"of":     MUL,
}

const (
	ERROR = "ERROR"
	NIL   = "NIL"
	EOF   = "EOF"

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
	NOT = "NOT"

	LPAREN  = "LPAREN"
	RPAREN  = "RPAREN"
	LSQUARE = "LSQUARE"
	RSQUARE = "RSQUARE"
	LBRACE  = "LBRACE"
	RBRACE  = "RBRACE"

	SEMICOLON = "SEMICOLON"
	COMMA     = "COMMA"
	ARROW     = "ARROW"
	TILDE     = "TILDE"
	DOT       = "DOT"

	IN     = "IN"
	DEFINE = "DEFINE"

	PROGRAM             = "PROGRAM"
	INT                 = "INT"
	STRING              = "STRING"
	FLOAT               = "FLOAT"
	IDENTIFIER          = "IDENTIFIER"
	UNIT                = "UNIT"
	BIN_OP              = "BIN_OP"
	UNARY_OP            = "UNARY_OP"
	FUNCTION_CALL       = "FUNCTION_CALL"
	ARRAY               = "ARRAY"
	ASSIGN              = "ASSIGN"
	FUNCTION_DEFENITION = "FUNCTION_DEFENITION"
	PERCENTAGE          = "PERCENTAGE"
	DICTIONARY          = "DICTIONARY"
	INDEX               = "INDEX"
)
