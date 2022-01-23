package evaluator

import (
	"errors"
	"math"
	"numerical/lexer"
	"numerical/parser"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	units "github.com/bcicen/go-units"
)

/* ----------------------------- Define Go Units ---------------------------- */

// Define any go units
func DefineUnits() {
	Week := units.NewUnit("week", "weeks")
	units.NewRatioConversion(Week, units.Day, 7.0)
}

/* --------------------------- Evaluator Functions -------------------------- */

// Generic Eval function
func Eval(node parser.Node, environment Environment) (Object, error) {
	switch n := node.(type) {
	case *parser.ProgramNode:
		program := Program{}
		for _, node := range n.Nodes {
			result, err := Eval(node, environment)
			if err != nil {
				return &Error{}, err
			}
			program.Objects = append(program.Objects, result)
		}
		return &program, nil

	case *parser.AssignNode:
		return evalAssign(n, environment)

	case *parser.ArrayNode:
		nodes, err := Eval(&n.Array, environment)
		if err != nil {
			return &Error{}, err
		}
		return &Array{nodes.(*Program).Objects}, nil

	case *parser.FunctionDefenitionNode:
		environment.Functions[n.Identifier] = n
		return &Nil{}, nil

	case *parser.UnitNode:
		value, err := Eval(n.Value, environment)
		return handleReturn(&Unit{Value: value.(Number).Inspect(), Unit: n.Unit}, err)

	case *parser.PercentageNode:
		value, err := Eval(n.Value, environment)
		return handleReturn(&Percentage{value.(Number).Inspect() / 100}, err)

	case *parser.IndexNode:
		return evalIndex(n, environment)

	case *parser.UnaryOpNode:
		return handleReturn(evalUnaryOp(n, environment))

	case *parser.BinOpNode:
		return handleReturn(evalBinaryOp(n, environment))

	case *parser.FunctionCallNode:
		return evalFunctionCall(n, environment)

	case *parser.IdentifierNode:
		return evalIdentifier(n, environment)

	case *parser.DictionaryNode:
		return evalDictionary(n, environment)

	case *parser.IntNode:
		return &Integer{Value: n.Value}, nil

	case *parser.FloatNode:
		return &Float{Value: n.Value}, nil

	case *parser.StringNode:
		return &String{Value: n.Value}, nil
	}

	return &Error{}, errors.New("EvaluationError: Undefined Node " + node.Type())
}

/* ------------------------- Extracted Eval Methods ------------------------- */

func evalAssign(n *parser.AssignNode, environment Environment) (Object, error) {
	value, err := Eval(n.Expression, environment)
	if err != nil {
		return &Error{}, err
	}
	environment.Variables[n.Identifier] = value
	return &Nil{}, nil
}

func evalIndex(n *parser.IndexNode, environment Environment) (Object, error) {
	array, err := Eval(n.Array, environment)
	if err != nil {
		return &Error{}, err
	}

	if array, ok := array.(*Array); ok {
		index, err := Eval(n.Index, environment)
		if err != nil {
			return &Error{}, err
		}
		if index, ok := index.(*Integer); ok {
			if int(index.Inspect()) >= len(array.Array) {
				return &Error{}, errors.New("IndexError: Index out of range")
			}
			return array.Array[int(index.Inspect())], nil
		} else {
			return &Error{}, errors.New("IndexError: Index is not type INT")
		}
	}
	return &Error{}, errors.New("IndexError: Array is not type ARRAY")
}

// Extracted <- eval class node (periodictable.hydrogen.mass)
func evalDictionary(n *parser.DictionaryNode, environment Environment) (Object, error) {
	container, err := Eval(n.Container, environment)
	if err != nil {
		return &Error{}, errors.New("ObjectError: Undefined identifier: " + container.String())
	}

	if container.Type() != lexer.DICTIONARY {
		return &Error{}, errors.New("ObjectError: Object returned not type DICTIONARY")
	}

	if value, ok := container.(*Dictionary).Dictionary[n.Field.Identifier]; ok {
		return value, nil
	}
	return &Error{}, errors.New("ObjectError: Undefined field referenced in object: " + n.Field.Identifier)
}

// Extracted <- eval identifier (periodic table, constants and variables)
func evalIdentifier(n *parser.IdentifierNode, environment Environment) (Object, error) {
	var keywords = map[string]func(Program, Environment) (Object, error){
		"prev":    prev,
		"history": history,
	}
	if f, ok := keywords[n.Identifier]; ok {
		result, err := f(Program{}, environment)
		if err != nil {
			return &Error{}, err
		}
		return result, nil
	}

	if value, ok := environment.Constants[n.Identifier]; ok {
		return value, nil
	}

	element, err := lookupElements(n.Identifier, environment.PeriodicTable)
	if err == nil {
		return element, nil
		// return formatFloat(element["atomic_mass"].(float64)), nil
	}

	if value, ok := environment.Variables[n.Identifier]; ok {
		return value, nil

	} else {
		if len(environment.Variables) < 1 {
			return &Error{}, errors.New("VarAccessError: Undefined variable identifier " + n.Identifier)
		}

		maxIdentifier := ""
		maxSimilarity := 0.0

		for variable := range environment.Variables {
			similarity := similarity(n.Identifier, variable)

			if similarity > maxSimilarity {
				maxSimilarity = similarity
				maxIdentifier = variable
			}
		}

		return environment.Variables[maxIdentifier], nil
	}
}

// Extracted <- function call (user defined or predefined)
func evalFunctionCall(n *parser.FunctionCallNode, environment Environment) (Object, error) {
	var functions = map[string]func(Program, Environment) (Object, error){
		"frac":    frac,
		"print":   print,
		"root":    root,
		"lookup":  lookup,
		"sin":     sin,
		"cos":     cos,
		"tan":     tan,
		"asin":    asin,
		"acos":    acos,
		"atan":    atan,
		"prev":    prev,
		"history": history,
		"sum":     sum,
	}

	if function, ok := functions[n.Identifier]; ok {
		params, err := Eval(&n.Parameters, environment)
		if err != nil {
			return &Error{}, err
		}
		if paramsProgram, ok := params.(*Program); ok {
			result, err := function(*paramsProgram, environment)
			if err != nil {
				return &Error{}, err
			}
			return result, nil
		}

	} else if function, ok := environment.Functions[n.Identifier]; ok {
		env := GenerateEnvironment()

		for i, node := range n.Parameters.Nodes {
			identifer := function.Parameters[i].(*parser.IdentifierNode).Identifier
			result, err := Eval(node, environment)
			if err != nil {
				return &Error{}, err
			}

			env.Variables[identifer] = result
		}

		result, err := Eval(&function.Consequence, env)
		if err != nil {
			return &Error{}, err
		}
		return result.(*Program).Objects[len(result.(*Program).Objects)-1], nil
	}

	return &Error{}, errors.New("FunctionCallError: Function with Identifer " + n.Identifier + " is not defined")
}

/* ---------------------------- Unary Operations ---------------------------- */

// Generic Unary Operation Function
func evalUnaryOp(node *parser.UnaryOpNode, environment Environment) (Object, error) {
	result, err := Eval(node.Right, environment)
	if err != nil {
		return &Error{}, err
	}

	switch node.Operation {
	case lexer.SUB:
		return evalUnarySub(result)
	case lexer.NOT:
		return evalUnaryNot(result), nil
	case lexer.TILDE:
		return evalUnaryRound(result)
	}

	return &Error{}, errors.New("UnaryOperationError: Unsupported " + node.Operation + " Operation")
}

/* --------------------------- Extracted Functions -------------------------- */

// Negate Operation -
func evalUnarySub(node Object) (Object, error) {
	switch n := node.(type) {
	case *Integer:
		return &Integer{Value: -n.Value}, nil
	case *Float:
		return &Float{Value: -n.Value}, nil
	}

	return &Error{}, errors.New("UnaryOperationError: Cannot negate type " + node.Type())
}

// Round to Integer ~
func evalUnaryRound(node Object) (Object, error) {
	switch n := node.(type) {
	case *Integer:
		return n, nil
	case *Float:
		return &Integer{Value: int(math.Round(n.Value))}, nil
	}

	return &Error{}, errors.New("RoundingError: cannout round type " + node.Type())
}

// Binary NOT operation !
func evalUnaryNot(node Object) *Integer {
	switch n := node.(type) {
	case *Integer:
		if n.Value == 0 {
			return &Integer{Value: 1}
		}
	case *String:
		if n.Value == "" {
			return &Integer{Value: 1}
		}
	}

	return &Integer{Value: 0}
}

/* ---------------------------- Binary Operations --------------------------- */

// Generic Binary Expression
func evalBinaryOp(node *parser.BinOpNode, environment Environment) (Object, error) {
	left, err := Eval(node.Left, environment)
	if err != nil {
		return &Error{}, err
	}

	shouldReturn, unitNode, err := handleInOperation(node, left)
	if shouldReturn {
		return unitNode, err
	}

	right, err := Eval(node.Right, environment)
	if err != nil {
		return &Error{}, err
	}

	if leftObj, ok := left.(Factor); ok {
		return leftObj.BinaryOperation(right, node.Operation)
	}

	return &Error{}, errors.New("BinaryOperationError: Unsupported Types: " + left.Type() + " " + node.Operation + " " + right.Type())
}

/* ---------------------------- Extracted Methods --------------------------- */

// Since IN requires Identifier, must be done before evaluated,
func handleInOperation(node *parser.BinOpNode, left Object) (bool, Object, error) {
	if node.Operation == lexer.IN && node.Right.Type() == lexer.IDENTIFIER {
		toIdentifier := node.Right.(*parser.IdentifierNode).Identifier

		if leftUnit, ok := left.(*Unit); ok {
			converted, err := convert(leftUnit.Inspect(), leftUnit.Unit, toIdentifier)
			return true, converted, err

		} else if leftUnit, ok := left.(Number); ok {
			converted, err := convert(leftUnit.Inspect(), toIdentifier, toIdentifier)
			return true, converted, err
		}

	} else if node.Operation == lexer.IN && node.Right.Type() != lexer.IDENTIFIER {
		return true, &Error{}, errors.New("ConversionError: IN cannot convert " + left.Type() + " and " + node.Right.Type())
	}
	return false, nil, nil
}

// Base Binary Operations With Floats
func binaryOperations(left float64, right float64, operation string) float64 {
	var result float64
	switch operation {
	case lexer.ADD:
		result = left + right
	case lexer.SUB:
		result = left - right
	case lexer.DIV:
		result = left / right
	case lexer.MUL:
		result = left * right
	case lexer.POW:
		result = math.Pow(left, right)
	case lexer.MOD:
		result = math.Mod(left, right)
	case lexer.EE:
		result = float64(boolToInt(left == right))
	case lexer.NE:
		result = float64(boolToInt(left != right))
	case lexer.LT:
		result = float64(boolToInt(left < right))
	case lexer.LTE:
		result = float64(boolToInt(left <= right))
	case lexer.GT:
		result = float64(boolToInt(left > right))
	case lexer.GTE:
		result = float64(boolToInt(left >= right))
	}
	return result
}

/* ---------------------------- Helper Functions ---------------------------- */

// Boolean Value to int
func boolToInt(value bool) int {
	if value {
		return 1
	} else {
		return 0
	}
}

// return INT if integer else format float to 5 d.p
func formatFloat(float float64) Number {
	if float64(int(float)) == float {
		return &Integer{Value: int(float)}
	}
	return &Float{Value: math.Round(float*100000) / 100000}
}

// Wrapper function for repeated code
func handleReturn(obj Object, err error) (Object, error) {
	if err != nil {
		return &Error{}, err
	}
	return obj, nil
}

// Convert UNITS
func convert(u float64, from string, to string) (unit *Unit, err error) {
	if from == to {
		return &Unit{Value: u, Unit: from}, nil
	}

	leftUnit, err := units.Find(from)
	if err != nil {
		return &Unit{}, errors.New("ConversionError: Unit " + from + " not defined")
	}

	rightUnit, err := units.Find(to)
	if err != nil {
		return &Unit{}, errors.New("ConversionError: Unit " + to + " not defined")
	}

	defer func() {
		if r := recover(); r != nil {
			err = nil
			unit = &Unit{u, to}
		}
	}()

	return &Unit{formatFloat(units.MustConvertFloat(u, leftUnit, rightUnit).Float()).Inspect(), to}, nil
}

// Find element identifier in periodic table and return element
func lookupElements(elementIdentifier string, periodicTable map[string]interface{}) (element Object, err error) {
	defer func() {
		if r := recover(); r != nil {
			element = &Dictionary{}
			err = errors.New("LookupError: Periodic Table is nil")
		}
	}()

	for _, element := range periodicTable["elements"].([]interface{}) {

		if element.(map[string]interface{})["symbol"].(string) == elementIdentifier {
			object, err := dictionaryFromMap(element.(map[string]interface{}))
			if err != nil {
				return &Error{}, err
			}
			return object, nil

		} else if strings.EqualFold(element.(map[string]interface{})["name"].(string), elementIdentifier) {
			object, err := dictionaryFromMap(element.(map[string]interface{}))
			if err != nil {
				return &Error{}, err
			}
			return object, nil
		}
	}

	return &Dictionary{}, errors.New("EvaluationError: Identifier undefined")
}

// Return similarity of strings, 50% consec and 50% Levenshtien
func similarity(a, b string) float64 {
	sequencer := strutil.Similarity(a, b, metrics.NewLevenshtein())

	i := 0
	for i < len(a) && i < len(b) {
		if a[i] != b[i] {
			break
		}
		i++
	}

	consecutiveCertainty := i / len(a)
	return (float64(consecutiveCertainty) * 0.5) + (sequencer * 0.5)
}

func dictionaryFromMap(m map[string]interface{}) (Object, error) {
	dictionaryObject := &Dictionary{Dictionary: make(map[string]Object)}
	for key, value := range m {
		switch value := value.(type) {
		case int:
			dictionaryObject.Dictionary[key] = formatFloat(float64(value))
		case float64:
			dictionaryObject.Dictionary[key] = formatFloat(value)
		case string:
			dictionaryObject.Dictionary[key] = &String{Value: value}
		case bool:
			dictionaryObject.Dictionary[key] = &Integer{Value: boolToInt(value)}

		}
	}
	return dictionaryObject, nil
}
