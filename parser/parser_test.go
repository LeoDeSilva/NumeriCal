package parser

import (
	"numerical/lexer"
	"reflect"
	"strconv"
	"testing"
)

func Test_preference(t *testing.T) {
	tests := []struct {
		tokenType string
		wantRes   int
	}{
		{lexer.ADD, 20},
		{lexer.MUL, 30},
		{lexer.EE, 10},
		{lexer.EOF, 0},
		{lexer.LPAREN, 0},
		{lexer.POW, 40},
		{lexer.ARROW, 5},
		{lexer.IN, 5},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := preference(tt.tokenType); got != tt.wantRes {
				t.Errorf("preference() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	operations := []string{"EE", "NE", "LT", "GT", "LTE", "GTE", "ADD", "SUB", "MUL", "DIV", "MOD", "POW", "IN", "ARROW"}
	tests := []struct {
		array   []string
		element string
		wantRes bool
	}{
		{operations, "IN", true},
		{operations, "ADD", true},
		{operations, "MOD", true},
		{operations, "EE", true},
		{operations, "INT", false},
		{operations, "NN", false},
		{operations, "", false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := contains(tt.array, tt.element); got != tt.wantRes {
				t.Errorf("contains() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		program string
		wantRes ProgramNode
		wantErr bool
	}{
		{"1+2", ProgramNode{lexer.PROGRAM_NODE, []Node{
			&BinOpNode{lexer.BIN_OP_NODE, &IntNode{lexer.INT_NODE, 1}, lexer.ADD, &IntNode{lexer.INT_NODE, 2}}},
		}, false},

		{"1-2", ProgramNode{lexer.PROGRAM_NODE, []Node{
			&BinOpNode{lexer.BIN_OP_NODE, &IntNode{lexer.INT_NODE, 1}, lexer.SUB, &IntNode{lexer.INT_NODE, 2}}},
		}, false},

		{"(1+2)*3", ProgramNode{lexer.PROGRAM_NODE, []Node{
			&BinOpNode{lexer.BIN_OP_NODE,
				&BinOpNode{lexer.BIN_OP_NODE, &IntNode{lexer.INT_NODE, 1}, lexer.ADD, &IntNode{lexer.INT_NODE, 2}},
				lexer.MUL,
				&IntNode{lexer.INT_NODE, 3}}},
		}, false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.NewLexer(tt.program)
			tokens, err := l.Lex()
			if err != nil {
				t.Errorf("Lexer.Lex() error = %v", err)
			}
			p := NewParser(tokens)
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantRes) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.wantRes)
			}
		})
	}
}
