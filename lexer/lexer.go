package lexer

import "fmt"

//TODO <- LexDouble: EQ, ARROW AND EE

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

func (l *Lexer) Lex() []Token {
	var tokens []Token
	for l.ch != 0 {
		tok := l.nextToken()
		tokens = append(tokens, tok)
		l.readChar()
		fmt.Println(tok)
	}
	tokens = append(tokens, NewToken(EOF, 0))
	return tokens
}

func (l *Lexer) nextToken() Token {
	var tok Token
	l.consumeWhitespace()
	switch l.ch {
	case '+':
		return NewToken(ADD, l.ch)
	case '-':
		return NewToken(SUB, l.ch)
	case '*':
		return NewToken(MUL, l.ch)
	case '/':
		return NewToken(DIV, l.ch)
	case '%':
		return NewToken(MOD, l.ch)
	case '^':
		return NewToken(POW, l.ch)
	case '(':
		return NewToken(LPAREN, l.ch)
	case ')':
		return NewToken(RPAREN, l.ch)
	case '[':
		return NewToken(LSQUARE, l.ch)
	case ']':
		return NewToken(RSQUARE, l.ch)
	case '{':
		return NewToken(LBRACE, l.ch)
	case '}':
		return NewToken(RBRACE, l.ch)
	case ';':
		return NewToken(SEMICOLON, l.ch)
	default:
		//TODO -> String, Identifier, Int, Float
		//FLOAT + INT ->  Create number string and see if contains '.'
		//STRING <- all between matching quotes
		//IDENTIFIER <- All letters + _
		if isLetter(l.ch) {

		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) consumeWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
