package evaluator

import (
	"fmt"
	"numerical/lexer"
	"numerical/parser"
)

type Environment struct {
	Variables     map[string]Object
	Functions     map[string]*parser.FunctionDefenitionNode
	PeriodicTable map[string]interface{}
}

type Object interface {
	Type() string
	String() string
}

type Number interface {
	Object
	Inspect() float64
}

type Nil struct{}

func (n *Nil) Type() string   { return lexer.NIL }
func (n *Nil) String() string { return "<nil>" }

type Error struct{}

func (e *Error) Type() string   { return lexer.ERROR }
func (e *Error) String() string { return "{ERROR}" }

type Program struct {
	Objects []Object
}

func (p *Program) Type() string { return lexer.PROGRAM_OBJ }
func (p *Program) String() string {
	repr := ""
	for i, node := range p.Objects {
		repr += node.String()
		if i != len(p.Objects)-1 {
			repr += "\n"
		}
	}
	return repr
}

type Unit struct {
	Value float64
	Unit  string
}

func (u *Unit) Type() string     { return lexer.UNIT_OBJ }
func (u *Unit) String() string   { return fmt.Sprintf("%v", u.Value) + " " + u.Unit }
func (u *Unit) Inspect() float64 { return u.Value }

type Integer struct {
	Value int
}

func (i *Integer) Type() string     { return lexer.INT_OBJ }
func (i *Integer) String() string   { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Inspect() float64 { return float64(i.Value) }

type Float struct {
	Value float64
}

func (f *Float) Type() string     { return lexer.FLOAT_OBJ }
func (f *Float) String() string   { return fmt.Sprintf("%v", f.Value) }
func (f *Float) Inspect() float64 { return f.Value }

type String struct {
	Value string
}

func (s *String) Type() string   { return lexer.STRING_OBJ }
func (s *String) String() string { return s.Value }
