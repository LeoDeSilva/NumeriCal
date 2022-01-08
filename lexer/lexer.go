package lexer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	program      string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(program string) *Lexer {
	l := &Lexer{program: program}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.program) {
		l.ch = 0
	} else {
		l.ch = l.program[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.program) {
		return 0
	} else {
		return l.program[l.readPosition]
	}
}

func (l *Lexer) Lex() ([]Token, error) {
	var tokens []Token
	for l.ch != 0 {
		tok, err := l.nextToken()
		if err != nil {
			return make([]Token, 0), err
		}
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, Token{EOF, ""})
	return tokens, nil
}

func (l *Lexer) nextToken() (Token, error) {
	var tok Token
	l.consumeWhitespace()
	switch l.ch {
	case '+':
		tok = NewToken(ADD, l.ch)
	case '-':
		tok = NewToken(SUB, l.ch)
	case '*':
		tok = NewToken(MUL, l.ch)
	case '/':
		tok = NewToken(DIV, l.ch)
	case '%':
		tok = NewToken(MOD, l.ch)
	case '^':
		tok = NewToken(POW, l.ch)
	case ',':
		tok = NewToken(COMMA, l.ch)
	case '.':
		tok = NewToken(DOT, l.ch)
	case '(':
		tok = NewToken(LPAREN, l.ch)
	case ')':
		tok = NewToken(RPAREN, l.ch)
	case '[':
		tok = NewToken(LSQUARE, l.ch)
	case ']':
		tok = NewToken(RSQUARE, l.ch)
	case '{':
		tok = NewToken(LBRACE, l.ch)
	case '}':
		tok = NewToken(RBRACE, l.ch)
	case ';':
		tok = NewToken(SEMICOLON, l.ch)
	case '~':
		tok = NewToken(TILDE, l.ch)
	case '=':
		tok = l.readDouble(EQ, map[byte]string{'=': EE, '>': ARROW})
	case '>':
		tok = l.readDouble(GT, map[byte]string{'=': GTE})
	case '<':
		tok = l.readDouble(LT, map[byte]string{'=': LTE})
	case '!':
		tok = l.readDouble(NOT, map[byte]string{'=': NE})
	case '\'', '"':
		var err error
		tok.Literal, err = l.readString(l.ch)
		tok.Type = STRING
		if err != nil {
			return Token{ERROR, ""}, err
		}
	default:
		if isLetter(l.ch) {
			var err error
			tok.Literal = l.readIdentifier()
			tok.Type, err = lookupIdentifier(tok.Literal)
			if err != nil {
				return Token{ERROR, ""}, err
			}
			return tok, nil
		}

		if '0' <= l.ch && l.ch <= '9' {
			var err error
			tok.Literal = l.readNumber()
			tok.Type, err = isFloat(tok.Literal)
			if err != nil {
				return Token{ERROR, ""}, err
			}
			return tok, nil
		}

		return NewToken(ERROR, l.ch), errors.New("SyntaxError: invalid character in lexer")
	}
	l.readChar()
	return tok, nil
}

func (l *Lexer) readDouble(primaryType string, secondaryTypes map[byte]string) Token {
	ch := l.ch
	if tokenType, ok := secondaryTypes[l.peekChar()]; ok {
		l.readChar()
		return Token{tokenType, string(ch) + string(l.ch)}
	} else {
		return Token{primaryType, string(ch)}
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.program[position:l.position]
}

func (l *Lexer) readString(ch byte) (string, error) {
	l.readChar()
	position := l.position
	for l.ch != ch {
		if l.ch == 0 {
			return "", errors.New("SyntaxError, EOL while scanning string literal")
		}
		l.readChar()
	}
	return l.program[position:l.position], nil
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.program[position:l.position]
}

func (l *Lexer) consumeWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func isFloat(number string) (string, error) {
	if number[len(number)-1] == '.' {
		return "", errors.New("SyntaxError: unexpected decimal point end of float")
	}

	if _, err := strconv.ParseFloat(number, 64); err != nil {
		fmt.Println(err)
		return "", errors.New("LexerError: non integer string passed to isFloat()")
	}

	count := strings.Count(number, ".")
	if count > 0 {
		if count > 1 {
			return "", errors.New("SyntaxError: > 1 decimal points in float")
		}
		return FLOAT, nil
	} else {
		return INT, nil
	}
}
