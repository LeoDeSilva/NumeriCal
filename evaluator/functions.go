package evaluator

import (
	"errors"
	"fmt"
	"math"
	"numerical/lexer"
	"strings"
)

func paramsToString(params Program) string {
	repr := ""
	for _, node := range params.Objects {
		repr += node.String()
	}
	return repr
}

func print(params Program, environment Environment) (Object, error) {
	fmt.Println(paramsToString(params))
	return &Nil{}, nil
}

func frac(params Program, environment Environment) (Object, error) {
	if len(params.Objects) < 2 {
		return &Error{}, errors.New("FracError: Expected parameters >= 2")
	}

	numerator := params.Objects[0]
	denomenator := params.Objects[1]

	if numerator, ok := numerator.(Number); ok {
		if denomenator, ok := denomenator.(Number); ok {
			result, err := evalNumberInfix(numerator, denomenator, lexer.DIV)
			if err != nil {
				return &Error{}, err
			}
			return result, nil
		}
	}
	return &Error{}, errors.New("BinaryOperationError: Cannot divide types, " + numerator.Type() + " and " + denomenator.Type())
}

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

func lookup(params Program, environment Environment) (Object, error) {
	if len(params.Objects) < 1 {
		return &Error{}, errors.New("LookupError: expected parameter length > 1")
	} else if params.Objects[0].Type() != lexer.STRING_OBJ {
		return &Error{}, errors.New("LookupError: expected type STRING or IDENTIFIER, not type " + params.Objects[0].Type())
	}

	element, err := lookupElements(params.Objects[0].(*String).Value, environment.PeriodicTable)
	if err != nil {
		return &Error{}, errors.New("LookupError: element " + params.Objects[0].(*String).Value + " does not exist")
	}

	for key, value := range element {
		if strings.Contains("name appearance atomic_mass category number period phase summary symbol shells", key) {
			fmt.Println(key, ":", value)
		}
	}

	return &Nil{}, nil

}
