package evaluator

import (
	"numerical/lexer"
	"numerical/parser"
	"reflect"
	"strconv"
	"testing"
)

func TestEval(t *testing.T) {
	defaultEnvironment := GenerateEnvironment()
	tests := []struct {
		program     string
		environment Environment
		want        []Object
		wantErr     bool
	}{
		//Data Types
		{program: "10", want: []Object{&Integer{10}}},
		{program: "10.2", want: []Object{&Float{10.2}}},
		{program: "10m", want: []Object{&Unit{10.0, "m"}}},
		{program: "10%", want: []Object{&Percentage{0.1}}},
		{program: "(10+2)%", want: []Object{&Percentage{0.12}}},
		{program: "'hello world'", want: []Object{&String{"hello world"}}},
		{program: "H", environment: defaultEnvironment, want: []Object{&Float{1.008}}},
		{program: "Hydrogen", environment: defaultEnvironment, want: []Object{&Float{1.008}}},
		{program: "rent", environment: modifyVariables(defaultEnvironment, "rent", &Integer{100}), want: []Object{&Integer{100}}},
		{program: "r", environment: modifyVariables(defaultEnvironment, "rent", &Integer{100}), want: []Object{&Integer{100}}},

		//Unary
		{program: "~10.2", want: []Object{&Integer{10}}},
		{program: "-10", want: []Object{&Integer{-10}}},
		{program: "!1", want: []Object{&Integer{0}}},
		{program: "!''", want: []Object{&Integer{1}}},
		{program: "~10", want: []Object{&Integer{10}}},
		{program: "-10.2", want: []Object{&Float{-10.2}}},

		//Binary
		{program: "1+1", want: []Object{&Integer{2}}},
		{program: "2-1", want: []Object{&Integer{1}}},
		{program: "2*4", want: []Object{&Integer{8}}},
		{program: "4/2", want: []Object{&Integer{2}}},
		{program: "4^2", want: []Object{&Integer{16}}},
		{program: "4%2", want: []Object{&Integer{0}}},
		{program: "4==2", want: []Object{&Integer{0}}},
		{program: "4!=2", want: []Object{&Integer{1}}},
		{program: "4>2", want: []Object{&Integer{1}}},
		{program: "4>=4", want: []Object{&Integer{1}}},
		{program: "4<2", want: []Object{&Integer{0}}},
		{program: "10<=10", want: []Object{&Integer{1}}},

		{program: "50 * 10%", want: []Object{&Integer{5}}},
		{program: "50 + 10%", want: []Object{&Float{50.1}}},
		{program: "10% + 10%", want: []Object{&Percentage{0.2}}},

		{program: "'hello ' + 'world'", want: []Object{&String{"hello world"}}},

		{program: "4km / 100m", want: []Object{&Unit{40, "m"}}},
		{program: "4 + 100m", want: []Object{&Unit{104, "m"}}},
		{program: "40m * 10", want: []Object{&Unit{400, "m"}}},
		{program: "4km / 100m in km", want: []Object{&Unit{0.04, "km"}}},
		{program: "40m * 100m in km", want: []Object{&Unit{4, "km"}}},

		{program: "100m in km", want: []Object{&Unit{0.1, "km"}}},

		//Functions
		{program: "frac(1,2)", want: []Object{&Float{0.5}}},
		{program: "root(8,3)", want: []Object{&Float{2}}},
		{program: "root(9)", want: []Object{&Integer{3}}},

		{program: "define f(x) => x^2", environment: defaultEnvironment, want: []Object{&Nil{}}},

		{program: "rent=100", environment: defaultEnvironment, want: []Object{&Nil{}}},
		{program: "rent=100; r", environment: defaultEnvironment, want: []Object{&Nil{}, &Integer{100}}},
		{program: "lookup('H')", environment: defaultEnvironment, want: []Object{&Nil{}}},
		{program: "print('hello world')", environment: defaultEnvironment, want: []Object{&Nil{}}},

		{
			program: "f(9)",
			environment: modifyFunctions(defaultEnvironment, "f", &parser.FunctionDefenitionNode{
				Identifier: "f",
				Parameters: []parser.Node{&parser.IdentifierNode{Identifier: "x"}},
				Consequence: parser.ProgramNode{Nodes: []parser.Node{&parser.BinOpNode{
					Left:      &parser.IdentifierNode{Identifier: "x"},
					Operation: lexer.POW,
					Right:     &parser.IntNode{Value: 2},
				}},
				}}),
			want: []Object{&Integer{Value: 81}},
		},

		//ERROR handling
		{program: "'hello world' * 'j'", wantErr: true},
		{program: "10+'hello world'", wantErr: true},
		{program: "f(8)", wantErr: true},
		{program: "10m in hours", wantErr: true},
		{program: "10msdf in hours", wantErr: true},
		{program: "10m in hosdfrs", wantErr: true},
		{program: "root()", wantErr: true},
		{program: "frac()", wantErr: true},
		{program: "lookup(10)", wantErr: true},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			l := lexer.NewLexer(tt.program)
			tokens, err := l.Lex()
			if err != nil {
				t.Errorf("Lexer.Lex() error = %v", err)
				return
			}
			p := parser.NewParser(tokens)
			ast, err := p.Parse()
			if err != nil {
				t.Errorf("Parser.Parse() error = %v", err)
				return
			}
			got, err := Eval(&ast, tt.environment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, &Program{tt.want}) && !tt.wantErr {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func modifyVariables(environment Environment, identifier string, value Object) Environment {
	environment.Variables[identifier] = value
	return environment
}

func modifyFunctions(environment Environment, identifier string, function *parser.FunctionDefenitionNode) Environment {
	environment.Functions[identifier] = function
	return environment
}
