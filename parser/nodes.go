package parser

import (
	"fmt"
	"numerical/lexer"
)

/* ---------------------------------- Nodes --------------------------------- */

// Wrapper for all nodes
type Node interface {
	Type() string
	String() string
}

/* ------------------------------ Wrapper nodes ----------------------------- */

// Error Wrapper
type ErrorNode struct{}

func (e *ErrorNode) Type() string   { return lexer.ERROR }
func (e *ErrorNode) String() string { return "{ERROR_NODE}" }

// []Nodes
type ProgramNode struct {
	Nodes []Node
}

func (p *ProgramNode) Type() string { return lexer.PROGRAM_NODE }
func (e *ProgramNode) String() string {
	repr := "["
	for i, node := range e.Nodes {
		repr += node.String()
		if i != len(e.Nodes)-1 {
			repr += ","
		}
	}
	return repr + "]"
}

/* ---------------------------- Expression Nodes ---------------------------- */

// define f(x) => x^2
type FunctionDefenitionNode struct {
	Identifier  string
	Parameters  []Node
	Consequence ProgramNode
}

func (f *FunctionDefenitionNode) Type() string { return lexer.FUNCTION_DEFENITION_NODE }
func (f *FunctionDefenitionNode) String() string {
	repr := "(" + f.Identifier + "("
	for _, node := range f.Parameters {
		repr += node.String()
	}
	repr += ") => "
	for _, node := range f.Consequence.Nodes {
		repr += node.String()
	}
	return repr + ")"
}

// rent = 10
type AssignNode struct {
	Identifier string
	Expression Node
}

func (a *AssignNode) Type() string { return lexer.ASSIGN_NODE }
func (a *AssignNode) String() string {
	return "(" + a.Identifier + "=" + a.Expression.String() + ")"
}

// 10+1, 10m in km
type BinOpNode struct {
	Left      Node
	Operation string
	Right     Node
}

func (b *BinOpNode) Type() string { return lexer.BIN_OP_NODE }
func (b *BinOpNode) String() string {
	return "(" + b.Left.String() + ":" + b.Operation + ":" + b.Right.String() + ")"
}

// -10, ~10.2
type UnaryOpNode struct {
	Operation string
	Right     Node
}

func (u *UnaryOpNode) Type() string   { return lexer.UNARY_OP_NODE }
func (u *UnaryOpNode) String() string { return "(" + u.Operation + ":" + u.Right.String() + ")" }

/* ------------------------------ Factor Nodes ------------------------------ */

// print(), frac(1,2)
type FunctionCallNode struct {
	Identifier string
	Parameters ProgramNode
}

func (f *FunctionCallNode) Type() string { return lexer.FUNCTION_CALL_NODE }
func (f *FunctionCallNode) String() string {
	repr := "(" + f.Identifier + "("
	for _, node := range f.Parameters.Nodes {
		repr += node.String()
	}
	repr += "))"
	return repr
}

// [1,2,3]
type ArrayNode struct {
	Nodes ProgramNode
}

func (a *ArrayNode) Type() string { return lexer.ARRAY_NODE }
func (a *ArrayNode) String() string {
	repr := "["
	for i, node := range a.Nodes.Nodes {
		repr += node.String()
		if i < len(a.Nodes.Nodes)-1 {
			repr += ","
		}
	}
	return repr + "]"
}

/* ------------------------------ Factor Nodes ------------------------------ */

// 10m, 10.2km
type UnitNode struct {
	Value Node
	Unit  string
}

func (u *UnitNode) Type() string   { return lexer.UNIT_NODE }
func (u *UnitNode) String() string { return u.Value.String() + u.Unit }

// 100% 10.2% 0.1%
type PercentageNode struct {
	Value Node
}

func (p *PercentageNode) Type() string   { return lexer.PERCENTAGE_NODE }
func (p *PercentageNode) String() string { return p.Value.String() + "%" }

// x, hello, rent
type IdentifierNode struct {
	Identifier string
}

func (i *IdentifierNode) Type() string   { return lexer.IDENTIFIER_NODE }
func (i *IdentifierNode) String() string { return i.Identifier }

// 10, 20
type IntNode struct {
	Value int
}

func (i *IntNode) Type() string   { return lexer.INT_NODE }
func (i *IntNode) String() string { return fmt.Sprintf("%d", i.Value) }

// 10.2, 10.3
type FloatNode struct {
	Value float64
}

func (f *FloatNode) Type() string   { return lexer.FLOAT_NODE }
func (f *FloatNode) String() string { return fmt.Sprintf("%v", f.Value) }

// "hello world"
type StringNode struct {
	Value string
}

func (s *StringNode) Type() string   { return lexer.STRING_NODE }
func (s *StringNode) String() string { return "\"" + s.Value + "\"" }
