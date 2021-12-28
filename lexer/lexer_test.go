package lexer

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func Test_lookupIdentifier(t *testing.T) {
	tests := []struct {
		args       string
		wantResult string
		wantErr    bool
	}{
		{args: "in", wantResult: IN},
		{args: "In", wantResult: IDENTIFIER},
		{args: "x", wantResult: IDENTIFIER},
		{args: "DEFINE", wantResult: IDENTIFIER},
		{args: "define", wantResult: DEFINE},
		{args: "Define", wantResult: IDENTIFIER},
		{args: "cuddles", wantResult: IDENTIFIER},
		{args: "", wantErr: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := lookupIdentifier((tt.args))
			if tt.wantErr {
				if err == nil {
					t.Errorf("lookupIdentifier() = Expected error, however unrecieved")
				}
			} else if err != nil {
				fmt.Println(err)
				t.Errorf("LookupIdentifier() = Unexpected error, %v", err)
			} else if got != tt.wantResult {
				t.Errorf("lookupIdentifier() = %v, want %v", got, tt.wantResult)
			}

		})
	}
}

func TestLexer_nextToken(t *testing.T) {
	tests := []struct {
		l          *Lexer
		wantResult Token
		wantErr    bool
	}{
		{l: NewLexer("+"), wantResult: Token{ADD, "+"}},
		{l: NewLexer("-"), wantResult: Token{SUB, "-"}},
		{l: NewLexer("*"), wantResult: Token{MUL, "*"}},
		{l: NewLexer("/"), wantResult: Token{DIV, "/"}},
		{l: NewLexer("("), wantResult: Token{LPAREN, "("}},
		{l: NewLexer("["), wantResult: Token{LSQUARE, "["}},
		{l: NewLexer("}"), wantResult: Token{RBRACE, "}"}},
		{l: NewLexer("^"), wantResult: Token{POW, "^"}},
		{l: NewLexer("hello"), wantResult: Token{IDENTIFIER, "hello"}},
		{l: NewLexer("'Hello World'"), wantResult: Token{STRING, "Hello World"}},
		{l: NewLexer("\"Hello World\""), wantResult: Token{STRING, "Hello World"}},
		{l: NewLexer("in"), wantResult: Token{IN, "in"}},
		{l: NewLexer("In"), wantResult: Token{IDENTIFIER, "In"}},
		{l: NewLexer("\"Hello World'"), wantResult: Token{ERROR, ""}, wantErr: true},
		{l: NewLexer("10"), wantResult: Token{INT, "10"}},
		{l: NewLexer("10.2"), wantResult: Token{FLOAT, "10.2"}},
		{l: NewLexer("10..2"), wantResult: Token{ERROR, ""}, wantErr: true},
		{l: NewLexer("10.2."), wantResult: Token{ERROR, ""}, wantErr: true},
		{l: NewLexer("."), wantResult: Token{ERROR, "."}, wantErr: true},
		{l: NewLexer(":"), wantResult: Token{ERROR, ":"}, wantErr: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := &Lexer{
				program:      tt.l.program,
				position:     tt.l.position,
				readPosition: tt.l.readPosition,
				ch:           tt.l.ch,
			}
			got, err := l.nextToken()
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("Lexer.nextToken() = %v, want %v", got, tt.wantResult)
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("nextToken() = Expected error, however unrecieved")
				}
			} else if err != nil {
				t.Errorf("LookupIdentifier() = Unexpected error, %v", err)
			}
		})
	}
}

func Test_isFloat(t *testing.T) {
	tests := []struct {
		number  string
		wantRes string
		wantErr bool
	}{
		{number: "10", wantRes: INT},
		{number: "10.2", wantRes: FLOAT},
		{number: "10.2.", wantErr: true},
		{number: "10.", wantErr: true},
		{number: "0", wantRes: INT},
		{number: "this", wantErr: true},
		{number: "10.h", wantErr: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := isFloat(tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("isFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRes {
				t.Errorf("isFloat() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func TestLexer_readIdentifier(t *testing.T) {
	tests := []struct {
		l    *Lexer
		want string
	}{
		{NewLexer("x"), "x"},
		{NewLexer("x."), "x"},
		{NewLexer("."), ""},
		{NewLexer(""), ""},
		{NewLexer("10"), ""},
		{NewLexer(".this"), ""},
		{NewLexer("hello world"), "hello"},
		{NewLexer("HelloWorld"), "HelloWorld"},
		{NewLexer("\"HelloWorld\""), ""},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := &Lexer{
				program:      tt.l.program,
				position:     tt.l.position,
				readPosition: tt.l.readPosition,
				ch:           tt.l.ch,
			}
			if got := l.readIdentifier(); got != tt.want {
				t.Errorf("Lexer.readIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readString(t *testing.T) {
	tests := []struct {
		l         *Lexer
		wantRes   string
		wantErr   bool
		terminate byte
	}{
		{l: NewLexer("\"HelloWorld\""), wantRes: "HelloWorld", terminate: '"'},
		{l: NewLexer("HelloWorld"), terminate: '"', wantErr: true},
		{l: NewLexer("'Hello WOrld\""), terminate: '\'', wantErr: true},
		{l: NewLexer("'Hello WOrld'"), terminate: '\'', wantRes: "Hello WOrld"},
		{l: NewLexer("''"), terminate: '\'', wantRes: ""},
		{l: NewLexer("\"\""), terminate: '"', wantRes: ""},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := &Lexer{
				program:      tt.l.program,
				position:     tt.l.position,
				readPosition: tt.l.readPosition,
				ch:           tt.l.ch,
			}
			got, err := l.readString(l.ch)
			if got != tt.wantRes {
				t.Errorf("Lexer.readString() = %v, want %v", got, tt.wantRes)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("readString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLexer_readNumber(t *testing.T) {
	tests := []struct {
		l       *Lexer
		wantRes string
	}{
		{NewLexer("10"), "10"},
		{NewLexer("10.2"), "10.2"},
		{NewLexer("10."), "10."},
		{NewLexer("hello"), ""},
		{NewLexer("10h"), "10"},
		{NewLexer("h10h"), ""},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := &Lexer{
				program:      tt.l.program,
				position:     tt.l.position,
				readPosition: tt.l.readPosition,
				ch:           tt.l.ch,
			}
			if got := l.readNumber(); got != tt.wantRes {
				t.Errorf("Lexer.readNumber() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		program string
		wantRes []Token
		wantErr bool
	}{
		{program: "10", wantRes: []Token{{INT, "10"}, {EOF, ""}}},
		{program: "10+2.5", wantRes: []Token{{INT, "10"}, {ADD, "+"}, {FLOAT, "2.5"}, {EOF, ""}}},
		{program: "10min in hours", wantRes: []Token{{INT, "10"}, {IDENTIFIER, "min"}, {IN, "in"}, {IDENTIFIER, "hours"}, {EOF, ""}}},    // in can also be =>
		{program: "10min => hours", wantRes: []Token{{INT, "10"}, {IDENTIFIER, "min"}, {ARROW, "=>"}, {IDENTIFIER, "hours"}, {EOF, ""}}}, // in can also be =>
		{program: ".10", wantRes: []Token{}, wantErr: true},
		{program: "'Helloworld", wantRes: []Token{}, wantErr: true},
		{program: "'Hello world'", wantRes: []Token{{STRING, "Hello world"}, {EOF, ""}}},
		{program: "sum(10,2,3)", wantRes: []Token{{IDENTIFIER, "sum"}, {LPAREN, "("}, {INT, "10"}, {COMMA, ","}, {INT, "2"}, {COMMA, ","}, {INT, "3"}, {RPAREN, ")"}, {EOF, ""}}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := NewLexer(tt.program)
			got, err := l.Lex()
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantRes) {
				t.Errorf("Lexer.Lex() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	tests := []struct {
		ch      byte
		wantRes bool
	}{
		{'a', true},
		{'z', true},
		{'h', true},
		{'_', true},
		{'0', false},
		{'$', false},
		{'"', false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := isLetter(tt.ch); got != tt.wantRes {
				t.Errorf("isLetter() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func Test_isDigit(t *testing.T) {
	tests := []struct {
		ch      byte
		wantRes bool
	}{
		{'0', true},
		{'9', true},
		{'a', false},
		{'.', true},
		{',', false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := isDigit(tt.ch); got != tt.wantRes {
				t.Errorf("isDigit() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func TestLexer_readDouble(t *testing.T) {
	tests := []struct {
		program        string
		primaryType    string
		secondaryTypes map[byte]string
		wantRes        Token
	}{
		{program: "==", primaryType: EQ, secondaryTypes: map[byte]string{'=': EE, '>': ARROW}, wantRes: Token{EE, "=="}},
		{program: "=", primaryType: EQ, secondaryTypes: map[byte]string{'=': EE, '>': ARROW}, wantRes: Token{EQ, "="}},
		{program: "=^", primaryType: EQ, secondaryTypes: map[byte]string{'=': EE, '>': ARROW}, wantRes: Token{EQ, "="}},
		{program: "!=", primaryType: NOT, secondaryTypes: map[byte]string{'=': NE}, wantRes: Token{NE, "!="}},
		{program: "!", primaryType: NOT, secondaryTypes: map[byte]string{'=': NE}, wantRes: Token{NOT, "!"}},
		{program: "=>", primaryType: EQ, secondaryTypes: map[byte]string{'>': ARROW, '=': EE}, wantRes: Token{ARROW, "=>"}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := NewLexer(tt.program)
			if got := l.readDouble(tt.primaryType, tt.secondaryTypes); !reflect.DeepEqual(got, tt.wantRes) {
				t.Errorf("Lexer.readDouble() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}
