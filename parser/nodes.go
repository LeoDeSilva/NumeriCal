package parser

import (
	"fmt"
	"strconv"
)

/* ---------------------------------- Nodes --------------------------------- */

type Node interface {
	Eval() Node
	String() string
}

type ErrorNode struct {
	Type string
}

func (e *ErrorNode) Eval() Node     { return e }
func (e *ErrorNode) String() string { return "{ERROR_NODE}" }

type ProgramNode struct {
	Type  string
	Nodes []Node
}

func (p *ProgramNode) Eval() Node { return p }
func (e *ProgramNode) String() string {
	repr := "{PROGRAM_NODE:"
	for _, node := range e.Nodes {
		repr += node.String()
	}
	return repr
}

type FunctionDefenitionNode struct {
	Type           string
	Identifier     string
	Configurations []Node
	Parameters     []Node
	Consequence    ProgramNode
}

func (f *FunctionDefenitionNode) Eval() Node { return f }
func (f *FunctionDefenitionNode) String() string {
	repr := "{DEFINE:" + f.Identifier + "("
	for _, node := range f.Configurations {
		repr += node.String()
	}
	repr += ")("
	for _, node := range f.Parameters {
		repr += node.String()
	}
	repr += ") => "
	for _, node := range f.Consequence.Nodes {
		repr += node.String()
	}
	return repr
}

type AssignNode struct {
	Type       string
	Identifier string
	Expression Node
}

func (a *AssignNode) Eval() Node { return a }
func (a *AssignNode) String() string {
	return "ASSIGN~{" + a.Identifier + "=" + a.Expression.String() + "}"
}

type BinOpNode struct {
	Type      string
	Left      Node
	Operation string
	Right     Node
}

func (b *BinOpNode) Eval() Node { return b }
func (b *BinOpNode) String() string {
	return "BINOP~{" + b.Left.String() + ":" + b.Operation + ":" + b.Right.String() + "}"
}

type FunctionCallNode struct {
	Type           string
	Identifier     string
	Configurations []Node
	Parameters     []Node
}

func (f *FunctionCallNode) Eval() Node { return f }
func (f *FunctionCallNode) String() string {
	repr := f.Identifier + "("
	for _, node := range f.Configurations {
		repr += node.String()
	}
	repr += ")("
	for _, node := range f.Parameters {
		repr += node.String()
	}
	repr += ")"
	return repr
}

type ArrayNode struct {
	Type  string
	Nodes []Node
}

func (a *ArrayNode) Eval() Node { return a }
func (a *ArrayNode) String() string {
	repr := "ARRAY~["
	for _, node := range a.Nodes {
		repr += node.String() + ","
	}
	return repr + "]"
}

type UnitNode struct {
	Type  string
	Value Node
	Unit  string
}

func (u *UnitNode) Eval() Node     { return u }
func (u *UnitNode) String() string { return "UNIT~" + u.Value.String() + u.Unit }

type IdentifierNode struct {
	Type       string
	Identifier string
}

func (i *IdentifierNode) Eval() Node     { return i }
func (i *IdentifierNode) String() string { return "IDENTIFIER~" + i.Identifier }

type IntNode struct {
	Type  string
	Value int
}

func (i *IntNode) Eval() Node     { return i }
func (i *IntNode) String() string { return "INT~" + strconv.Itoa(i.Value) }

type FloatNode struct {
	Type  string
	Value float64
}

func (f *FloatNode) Eval() Node     { return f }
func (f *FloatNode) String() string { return "FLOAT~" + fmt.Sprintf("%v", f.Value) }

type StringNode struct {
	Type  string
	Value string
}

func (s *StringNode) Eval() Node     { return s }
func (s *StringNode) String() string { return "STRING~" + "\"" + s.Value + "\"" }
