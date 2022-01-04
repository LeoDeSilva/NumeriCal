package evaluator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"numerical/lexer"
	"numerical/parser"
)

/* ------------------------ Environment and Functions ----------------------- */

// Contains Variables, Functions and Others
type Environment struct {
	Variables     map[string]Object
	Constants     map[string]Object
	Functions     map[string]*parser.FunctionDefenitionNode
	PeriodicTable map[string]interface{}
}

// Create constants for environment.Constants
func GenerateConstants() map[string]Object {
	return map[string]Object{
		"PI": formatFloat(3.141592),
		"E":  formatFloat(2.718281),
	}
}

func GenerateEnvironment() Environment {
	periodicTable, _ := ioutil.ReadFile("/Users/ldesilva/Documents/Personal/Coding/Golang/NumeriCal/evaluator/periodicTable.json")

	var periodicTableJson map[string]interface{}
	json.Unmarshal([]byte(periodicTable), &periodicTableJson)

	return Environment{
		Variables:     make(map[string]Object),
		Functions:     make(map[string]*parser.FunctionDefenitionNode),
		PeriodicTable: periodicTableJson,
		Constants:     GenerateConstants(),
	}
}

/* --------------------------------- Objects -------------------------------- */

// Wrapper for all nodes
type Object interface {
	Type() string
	String() string
}

/* ----------------------------- Wrapper Objects ---------------------------- */

// Wrapper for nil values returned
type Nil struct{}

func (n *Nil) Type() string   { return lexer.NIL }
func (n *Nil) String() string { return "<nil>" }

// Wrapper for all errors
type Error struct{}

func (e *Error) Type() string   { return lexer.ERROR }
func (e *Error) String() string { return "{ERROR}" }

/* ------------------------------ Node Objects ------------------------------ */

// Entire program [10+10, 2+2, print("hello world")]
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

// "Hello World"
type String struct {
	Value string
}

func (s *String) Type() string   { return lexer.STRING_OBJ }
func (s *String) String() string { return s.Value }

/* ----------------------------- Number Objects ----------------------------- */

// Wrapper for floats e.g. (UNIT, INT, PERCENTAGE)
type Number interface {
	Object
	Inspect() float64
}

// 10km etc
type Unit struct {
	Value float64
	Unit  string
}

func (u *Unit) Type() string     { return lexer.UNIT_OBJ }
func (u *Unit) String() string   { return fmt.Sprintf("%v", u.Value) + " " + u.Unit }
func (u *Unit) Inspect() float64 { return u.Value }

// 10, 12 etc
type Integer struct {
	Value int
}

func (i *Integer) Type() string     { return lexer.INT_OBJ }
func (i *Integer) String() string   { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Inspect() float64 { return float64(i.Value) }

// 10.2, 20.4 etc
type Float struct {
	Value float64
}

func (f *Float) Type() string     { return lexer.FLOAT_OBJ }
func (f *Float) String() string   { return fmt.Sprintf("%v", f.Value) }
func (f *Float) Inspect() float64 { return f.Value }
