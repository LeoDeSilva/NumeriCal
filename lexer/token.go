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
}

const (
	ERROR = "ERROR"
	EOF   = "EOF"

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

	IN     = "IN"
	DEFINE = "DEFINE"

	PROGRAM_NODE             = "PROGRAM_NODE"
	IDENTIFIER_NODE          = "IDENTIFIER_NODE"
	INT_NODE                 = "INT_NODE"
	STRING_NODE              = "STRING_NODE"
	FLOAT_NODE               = "FLOAT_NODE"
	UNIT_NODE                = "UNIT_NODE"
	BIN_OP_NODE              = "BIN_OP_NODE"
	UNARY_OP_NODE            = "UNARY_OP_NODE"
	FUNCTION_CALL_NODE       = "FUNCTION_CALL_NODE"
	ARRAY_NODE               = "ARRAY_NODE"
	ASSIGN_NODE              = "ASSIGN_NODE"
	FUNCTION_DEFENITION_NODE = "FUNCTION_DEFENITION_NODE"
)
