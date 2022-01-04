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
		wantRes []Node
		wantErr bool
	}{
		//Factors
		{"i", []Node{
			&IdentifierNode{"i"},
		}, false},
		{"10", []Node{&IntNode{10}}, false},
		{"10.2", []Node{&FloatNode{10.2}}, false},
		{"10m", []Node{&UnitNode{&IntNode{10}, "m"}}, false},
		{"\"hello world\"", []Node{&StringNode{"hello world"}}, false},
		{"[1,2,3]", []Node{&ArrayNode{ProgramNode{[]Node{&IntNode{1}, &IntNode{2}, &IntNode{3}}}}}, false},
		{"10%", []Node{&PercentageNode{&IntNode{10}}}, false},
		{"rent%", []Node{&PercentageNode{&IdentifierNode{"rent"}}}, false},
		{"(10+10)%", []Node{&PercentageNode{&BinOpNode{&IntNode{10}, lexer.ADD, &IntNode{10}}}}, false},

		//Function Calls
		{"print()", []Node{
			&FunctionCallNode{Identifier: "print", Parameters: ProgramNode{[]Node{}}},
		}, false},
		{"frac(1,2)", []Node{
			&FunctionCallNode{Identifier: "frac", Parameters: ProgramNode{[]Node{&IntNode{1}, &IntNode{2}}}},
		}, false},
		{"define f(x) => x^2", []Node{
			&FunctionDefenitionNode{
				Identifier: "f",
				Parameters: []Node{&IdentifierNode{Identifier: "x"}},
				Consequence: ProgramNode{[]Node{
					&BinOpNode{
						&IdentifierNode{Identifier: "x"},
						lexer.POW,
						&IntNode{2},
					}}}},
		}, false},

		//Unary Expressions
		{"~10", []Node{&UnaryOpNode{lexer.TILDE, &IntNode{10}}}, false},
		{"!1", []Node{&UnaryOpNode{lexer.NOT, &IntNode{1}}}, false},
		{"-10", []Node{&UnaryOpNode{lexer.SUB, &IntNode{10}}}, false},
		{"-(1+2)", []Node{&UnaryOpNode{lexer.SUB, &BinOpNode{&IntNode{1}, lexer.ADD, &IntNode{2}}}}, false},

		//Binary Expressions
		{"1+2;3+4", []Node{
			&BinOpNode{&IntNode{1}, lexer.ADD, &IntNode{2}},
			&BinOpNode{&IntNode{3}, lexer.ADD, &IntNode{4}},
		}, false},
		{"1+2", []Node{
			&BinOpNode{&IntNode{1}, lexer.ADD, &IntNode{2}},
		}, false},

		{"1-2", []Node{
			&BinOpNode{&IntNode{1}, lexer.SUB, &IntNode{2}},
		}, false},

		{"(1+2)*3", []Node{
			&BinOpNode{
				&BinOpNode{&IntNode{1}, lexer.ADD, &IntNode{2}},
				lexer.MUL,
				&IntNode{3}},
		}, false},
		{"rent=100", []Node{
			&AssignNode{"rent", &IntNode{100}},
		}, false},

		//Errors
		{"(", []Node{}, true},
		{"rent =", []Node{}, true},
		{"+10", []Node{}, true},
		{"%", []Node{}, true},
		{"%10", []Node{}, true},
		{"%rent", []Node{}, true},
		{"'rent'%", []Node{}, true},
		{"%(10+10)", []Node{}, true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.NewLexer(tt.program)
			tokens, err := l.Lex()
			if err != nil {
				t.Errorf("Lexer.Lex() error = %v", err)
				return
			}
			p := NewParser(tokens)
			got, err := p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() %v error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, ProgramNode{tt.wantRes}) && !tt.wantErr {
				t.Errorf("Parser.Parse() = %v, want %v", got, ProgramNode{tt.wantRes})
			}
		})
	}
}
