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
}

// Create constants for environment.Constants
func GenerateConstants() map[string]Object {
	return map[string]Object{
		"PI": formatFloat(3.141592),
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

/* ----------------------------- Factor Objects ----------------------------- */

// Wrapper for all objects that support binary operations
type Factor interface {
	Object
	BinaryOperation(Object, string) (Object, error)
}

// "Hello World"
type String struct {
	Value string
}

func (s *String) Type() string   { return lexer.STRING_OBJ }
func (s *String) String() string { return s.Value }
func (s *String) BinaryOperation(right Object, operation string) (Object, error) {
	if stringRight, ok := right.(*String); ok {
		switch operation {
		case lexer.ADD:
			return &String{s.Value + stringRight.Value}, nil
		default:
			return &Error{}, errors.New("BinaryOperationError: Undefined Operation " + operation + " on type STRING_OBJ")
		}
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types STRING_OBJ and " + right.Type())
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

func (u *Unit) Type() string     { return lexer.UNIT_OBJ }
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

	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types UNIT_OBJ and " + right.Type())
}

// 10%, 15%
type Percentage struct {
	Value float64
}

func (p *Percentage) Type() string     { return lexer.PERCENTAGE_OBJ }
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
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types PERCENTAGE_OBJ and " + right.Type())
}

// 10, 12 etc
type Integer struct {
	Value int
}

func (i *Integer) Type() string     { return lexer.INT_OBJ }
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
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types INT_OBJ and " + right.Type())
}

// 10.2, 20.4 etc
type Float struct {
	Value float64
}

func (f *Float) Type() string     { return lexer.FLOAT_OBJ }
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
	return &Error{}, errors.New("BinaryOperationError: Cannot operate on types FLOAT_OBJ and " + right.Type())
}
