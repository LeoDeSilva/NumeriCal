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
		wantError  bool
	}{
		{args: "in", wantResult: IN},
		{args: "In", wantResult: IDENTIFIER},
		{args: "x", wantResult: IDENTIFIER},
		{args: "DEFINE", wantResult: IDENTIFIER},
		{args: "define", wantResult: DEFINE},
		{args: "Define", wantResult: IDENTIFIER},
		{args: "cuddles", wantResult: IDENTIFIER},
		{args: "", wantError: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := lookupIdentifier((tt.args))
			if tt.wantError {
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
		l    *Lexer
		want Token
	}{
		{NewLexer("+"), Token{ADD, "+"}},
		{NewLexer("-"), Token{SUB, "-"}},
		{NewLexer("*"), Token{MUL, "*"}},
		{NewLexer("/"), Token{DIV, "/"}},
		{NewLexer("("), Token{LPAREN, "("}},
		{NewLexer("["), Token{LSQUARE, "["}},
		{NewLexer("}"), Token{RBRACE, "}"}},
		{NewLexer("^"), Token{POW, "^"}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := &Lexer{
				program:      tt.l.program,
				position:     tt.l.position,
				readPosition: tt.l.readPosition,
				ch:           tt.l.ch,
			}
			if got := l.nextToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexer.nextToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
