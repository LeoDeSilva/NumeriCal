package evaluator

import (
	"encoding/json"
	"errors"
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
	History       Array
}

// Create constants for environment.Constants
func GenerateConstants() map[string]Object {
	return map[string]Object{
		"Pi": formatFloat(3.141592),
		"E":  formatFloat(2.718281),
	}
}

// Generates the environment and loads in the periodic table
func GenerateEnvironment() Environment {
	periodicTable, _ := ioutil.ReadFile("/Users/ldesilva/Documents/Personal/Coding/Golang/NumeriCal/evaluator/periodicTable.json")

	var periodicTableJson map[string]interface{}
	json.Unmarshal([]byte(periodicTable), &periodicTableJson)

	return Environment{
		Variables:     make(map[string]Object),
		Functions:     make(map[string]*parser.FunctionDefenitionNode),
		PeriodicTable: periodicTableJson,
		Constants:     GenerateConstants(),
		History:       Array{make([]Object, 0)},
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
func (n *Nil) String() string { return "" }

// Wrapper for all errors
type Error struct{}

func (e *Error) Type() string   { return lexer.ERROR }
func (e *Error) String() string { return "{ERROR}" }

/* ------------------------------ Node Objects ------------------------------ */

// Entire program [10+10, 2+2, print("hello world")]
type Program struct {
	Objects []Object
}

func (p *Program) Type() string { return lexer.PROGRAM }
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

type Dictionary struct {
	Dictionary map[string]Object
}

func (d *Dictionary) Type() string   { return lexer.DICTIONARY }
func (d *Dictionary) String() string { return fmt.Sprintf("%T", d.Dictionary) }

/* ----------------------------- Factor Objects ----------------------------- */

// Wrapper for all objects that support binary operations
type Factor interface {
	Object
	BinaryOperation(Object, string) (Object, error)
}

// Array type elements are of dynamic type
type Array struct {
	Array []Object
}

func (a *Array) Type() string { return lexer.ARRAY }
func (a *Array) String() string {
	repr := "["
	for i, node := range a.Array {
		repr += node.String()
		if i < len(a.Array)-1 {
			repr += ","
		}
	}
	return repr + "]"
}

func (a *Array) BinaryOperation(right Object, operation string) (Object, error) {
	if operation != lexer.ADD {
		return &Error{}, errors.New("BinaryOperationError: Unsupported Operation on Array, " + operation)
	}

	if r, ok := right.(*Array); ok {
		return &Array{append(a.Array, r.Array...)}, nil
	}

	return &Array{append(a.Array, right)}, nil
}

// "Hello World"
type String struct {
	Value string
}

func (s *String) Type() string   { return lexer.STRING }
func (s *String) String() string { return s.Value }
func (s *String) BinaryOperation(right Object, operation string) (Object, error) {
	if stringRight, ok := right.(*String); ok {
		switch operation {
		case lexer.ADD:
			return &String{s.Value + stringRight.Value}, nil
		default:
			return &Error{}, errors.New("BinaryOperationError: Undefined Operation " + operation + " on type STRING")
		}
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types STRING and " + right.Type())
}

/* ----------------------------- Number Objects ----------------------------- */

// Wrapper for floats e.g. (UNIT, INT, PERCENTAGE)
type Number interface {
	Factor
	Inspect() float64
}

// 10km etc
type Unit struct {
	Value float64
	Unit  string
}

func (u *Unit) Type() string     { return lexer.UNIT }
func (u *Unit) String() string   { return fmt.Sprintf("%v", u.Value) + " " + u.Unit }
func (u *Unit) Inspect() float64 { return u.Value }
func (u *Unit) BinaryOperation(right Object, operation string) (Object, error) {
	switch right := right.(type) {
	case *Unit:
		convertedLeft, err := convert(u.Inspect(), u.Unit, right.Unit)
		if err != nil {
			return &Error{}, err
		}

		return &Unit{
			Value: binaryOperations(convertedLeft.Inspect(), right.Inspect(), operation),
			Unit:  right.Unit,
		}, nil

	case Number:
		return &Unit{
			Value: binaryOperations(u.Inspect(), right.Inspect(), operation),
			Unit:  u.Unit,
		}, nil
	}

	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types UNIT and " + right.Type())
}

// 10%, 15%
type Percentage struct {
	Value float64
}

func (p *Percentage) Type() string     { return lexer.PERCENTAGE }
func (p *Percentage) String() string   { return fmt.Sprintf("%d", int(p.Value*100)) + "%" }
func (p *Percentage) Inspect() float64 { return p.Value }
func (p *Percentage) BinaryOperation(right Object, operation string) (Object, error) {
	switch right := right.(type) {
	case *Percentage:
		return &Percentage{binaryOperations(p.Inspect(), right.Inspect(), operation)}, nil
	case Number:
		result := formatFloat(binaryOperations(p.Inspect(), right.Inspect(), operation))
		if right, ok := right.(*Unit); ok {
			return &Unit{Value: result.Inspect(), Unit: right.Unit}, nil
		}
		return result, nil

	}
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types PERCENTAGE and " + right.Type())
}

// 10, 12 etc
type Integer struct {
	Value int
}

func (i *Integer) Type() string     { return lexer.INT }
func (i *Integer) String() string   { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Inspect() float64 { return float64(i.Value) }
func (i *Integer) BinaryOperation(right Object, operation string) (Object, error) {
	switch right := right.(type) {
	case Number:
		result := formatFloat(binaryOperations(i.Inspect(), right.Inspect(), operation))
		if right, ok := right.(*Unit); ok {
			return &Unit{Value: result.Inspect(), Unit: right.Unit}, nil
		}
		return result, nil
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types INT and " + right.Type())
}

// 10.2, 20.4 etc
type Float struct {
	Value float64
}

func (f *Float) Type() string     { return lexer.FLOAT }
func (f *Float) String() string   { return fmt.Sprintf("%v", f.Value) }
func (f *Float) Inspect() float64 { return f.Value }
func (f *Float) BinaryOperation(right Object, operation string) (Object, error) {
	switch right := right.(type) {
	case Number:
		result := formatFloat(binaryOperations(f.Inspect(), right.Inspect(), operation))
		if right, ok := right.(*Unit); ok {
			return &Unit{Value: result.Inspect(), Unit: right.Unit}, nil
		}
		return result, nil
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types FLOAT and " + right.Type())
}
