package parser

import (
	"fmt"
	"numerical/lexer"
)

/* ---------------------------------- Nodes --------------------------------- */

type Node interface {
	Type() string
	String() string
}

type ErrorNode struct{}

func (e *ErrorNode) Type() string   { return lexer.ERROR }
func (e *ErrorNode) String() string { return "{ERROR_NODE}" }

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

type AssignNode struct {
	Identifier string
	Expression Node
}

func (a *AssignNode) Type() string { return lexer.ASSIGN_NODE }
func (a *AssignNode) String() string {
	return "(" + a.Identifier + "=" + a.Expression.String() + ")"
}

type BinOpNode struct {
	Left      Node
	Operation string
	Right     Node
}

func (b *BinOpNode) Type() string { return lexer.BIN_OP_NODE }
func (b *BinOpNode) String() string {
	return "(" + b.Left.String() + ":" + b.Operation + ":" + b.Right.String() + ")"
}

type UnaryOpNode struct {
	Operation string
	Right     Node
}

func (u *UnaryOpNode) Type() string   { return lexer.UNARY_OP_NODE }
func (u *UnaryOpNode) String() string { return "(" + u.Operation + ":" + u.Right.String() + ")" }

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

type ArrayNode struct {
	Nodes ProgramNode
}

func (a *ArrayNode) Type() string { return lexer.ARRAY_NODE }
func (a *ArrayNode) String() string {
	repr := "["
	for _, node := range a.Nodes.Nodes {
		repr += node.String() + ","
	}
	return repr + "]"
}

type UnitNode struct {
	Value Node
	Unit  string
}

func (u *UnitNode) Type() string   { return lexer.UNIT_NODE }
func (u *UnitNode) String() string { return u.Value.String() + u.Unit }

type IdentifierNode struct {
	Identifier string
}

func (i *IdentifierNode) Type() string   { return lexer.IDENTIFIER_NODE }
func (i *IdentifierNode) String() string { return i.Identifier }

type IntNode struct {
	Value int
}

func (i *IntNode) Type() string   { return lexer.INT_NODE }
func (i *IntNode) String() string { return fmt.Sprintf("%d", i.Value) }

type FloatNode struct {
	Value float64
}

func (f *FloatNode) Type() string   { return lexer.FLOAT_NODE }
func (f *FloatNode) String() string { return fmt.Sprintf("%v", f.Value) }

type StringNode struct {
	Value string
}

func (s *StringNode) Type() string   { return lexer.STRING_NODE }
func (s *StringNode) String() string { return "\"" + s.Value + "\"" }
