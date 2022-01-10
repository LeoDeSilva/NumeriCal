package evaluator

import (
	"errors"
	"fmt"
	"math"
	"numerical/lexer"
	"strings"
)

/* -------------------------------- Functions ------------------------------- */

// Print a string representation of the parameters
func print(params Program, environment Environment) (Object, error) {
	fmt.Println(paramsToString(params))
	return &Nil{}, nil
}

// Just divide 2 numbers frac(NUMERATOR, DENOMENATOR)
func frac(params Program, environment Environment) (Object, error) {
	if len(params.Objects) < 2 {
		return &Error{}, errors.New("FracError: Expected parameters >= 2")
	}

	numerator, okN := params.Objects[0].(Factor)
	denomenator, okD := params.Objects[1].(Factor)

	if okN && okD {
		result, err := numerator.BinaryOperation(denomenator, lexer.DIV)
		if err != nil {
			return &Error{}, err
		}
		return result, nil
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot divide types, " + numerator.Type() + " and " + denomenator.Type())
}

// Root a number (default root is 2) root(BASE, ROOT)
func root(params Program, environment Environment) (Object, error) {
	root := 2.0
	if len(params.Objects) < 1 {
		return &Error{}, errors.New("RootError: Expected parameters > 1")
	} else if len(params.Objects) >= 2 {
		if exponent, ok := params.Objects[1].(Number); ok {
			root = exponent.Inspect()
		} else {
			return &Error{}, errors.New("RootError: cannot use type " + params.Objects[1].Type() + " as exponent")
		}
	}
	if base, ok := params.Objects[0].(Number); ok {
		return formatFloat(math.Pow(base.Inspect(), 1/root)), nil
	} else {
		return &Error{}, errors.New("RootError: cannot raise type " + params.Objects[0].Type() + " to a power")
	}
}

// Lookup element in periodic table and return JSON lookup(ELEMENT)
// TODO: create element objects and return that instead
func lookup(params Program, environment Environment) (Object, error) {
	if len(params.Objects) < 1 {
		return &Error{}, errors.New("LookupError: expected parameter length > 1")
	} else if params.Objects[0].Type() != lexer.STRING {
		return &Error{}, errors.New("LookupError: expected type STRING or IDENTIFIER, not type " + params.Objects[0].Type())
	}

	element, err := lookupElements(params.Objects[0].(*String).Value, environment.PeriodicTable)
	if err != nil {
		return &Error{}, errors.New("LookupError: element " + params.Objects[0].(*String).Value + " does not exist")
	}

	for key, value := range element.(*Dictionary).Dictionary {
		if strings.Contains("name appearance atomic_mass category number period phase summary symbol shells", key) {
			fmt.Println(key, ":", value)
		}
	}

	return &Nil{}, nil

}

/* ------------------------- Trigonometry Functions ------------------------- */

// Handle calling trig functions that input in radians
func callTrig(params Program, function func(float64) float64) (Object, error) {
	if len(params.Objects) < 1 {
		return &Error{}, errors.New("FunctionErro: Expected parameter length > 1")
	} else if operand, ok := params.Objects[0].(Number); ok {
		return formatFloat(function(operand.Inspect() * (math.Pi / 180))), nil
	}

	return &Error{}, errors.New("FunctionError: Expected parameter type NUMBER")
}

// Handle calling trig functions that output in radians
func callReverseTrig(params Program, function func(float64) float64) (Object, error) {
	if len(params.Objects) < 1 {
		return &Error{}, errors.New("FunctionErro: Expected parameter length > 1")
	} else if operand, ok := params.Objects[0].(Number); ok {
		return formatFloat(function(operand.Inspect()) * (180 / math.Pi)), nil
	}

	return &Error{}, errors.New("FunctionError: Expected parameter type NUMBER")
}

func sin(params Program, environment Environment) (Object, error) {
	return callTrig(params, math.Sin)
}

func cos(params Program, environment Environment) (Object, error) {
	return callTrig(params, math.Cos)
}

func tan(params Program, environment Environment) (Object, error) {
	return callTrig(params, math.Tan)
}

func asin(params Program, environment Environment) (Object, error) {
	return callReverseTrig(params, math.Asin)
}

func acos(params Program, environment Environment) (Object, error) {
	return callReverseTrig(params, math.Acos)
}

func atan(params Program, environment Environment) (Object, error) {
	return callReverseTrig(params, math.Atan)
}

/* ---------------------------- Helper Functions ---------------------------- */

// Join parameters into string
func paramsToString(params Program) string {
	repr := ""
	for _, node := range params.Objects {
		repr += node.String()
	}
	return repr
}
