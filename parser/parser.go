package parser

import (
	"errors"
	"fmt"
	"numerical/lexer"
	"strconv"
)

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

type BinOpNode struct {
	Type      string
	Left      Node
	Operation string
	Right     Node
}

func (b *BinOpNode) Eval() Node { return b }
func (b *BinOpNode) String() string {
	return "{" + b.Left.String() + ":" + b.Operation + ":" + b.Right.String() + "}"
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
	repr := "["
	for _, node := range a.Nodes {
		repr += node.String()
	}
	return repr + "]"
}

type UnitNode struct {
	Type  string
	Value Node
	Unit  string
}

func (u *UnitNode) Eval() Node     { return u }
func (u *UnitNode) String() string { return u.Value.String() + u.Unit }

type IdentifierNode struct {
	Type       string
	Identifier string
}

func (i *IdentifierNode) Eval() Node     { return i }
func (i *IdentifierNode) String() string { return i.Identifier }

type IntNode struct {
	Type  string
	Value int
}

func (i *IntNode) Eval() Node     { return i }
func (i *IntNode) String() string { return strconv.Itoa(i.Value) }

type FloatNode struct {
	Type  string
	Value float64
}

func (f *FloatNode) Eval() Node     { return f }
func (f *FloatNode) String() string { return fmt.Sprintf("%v", f.Value) }

type StringNode struct {
	Type  string
	Value string
}

func (s *StringNode) Eval() Node     { return s }
func (s *StringNode) String() string { return "\"" + s.Value + "\"" }

// ===========PARSER==========

type Parser struct {
	tokens   []lexer.Token
	token    lexer.Token
	position int
}

func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{tokens: tokens, token: tokens[0], position: 0}
	return p
}

func (p *Parser) advance() lexer.Token {
	p.position++
	p.token = p.tokens[p.position]
	return p.tokens[p.position]
}

func (p *Parser) peekToken() lexer.Token {
	if p.position+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF, Literal: ""}
	}
	return p.tokens[p.position+1]
}

func (p *Parser) Parse() (ProgramNode, error) {
	ast := ProgramNode{Type: lexer.PROGRAM_NODE, Nodes: make([]Node, 0)}
	for p.token.Type != lexer.EOF {
		node, err := p.parseExpression()
		if err != nil {
			return ProgramNode{Type: lexer.PROGRAM_NODE}, err
		}
		ast.Nodes = append(ast.Nodes, node)
	}
	return ast, nil
}

func (p *Parser) parseExpression() (Node, error) {
	expr, err := p.parseExpr(0)
	if err != nil {
		return &ErrorNode{lexer.ERROR}, err
	}
	return expr, nil
}

func (p *Parser) parseExpr(rbp int) (Node, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return &ErrorNode{lexer.ERROR}, err
	}
	peekRbp := preference(p.token.Type)
	for p.peekToken().Type != lexer.EOF && peekRbp >= rbp {
		left, err = p.parseInfix(left, p.token.Type)
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		peekRbp = preference(p.token.Type)
	}
	return left, nil
}

func (p *Parser) parsePrefix() (Node, error) {
	switch p.token.Type {
	case lexer.LPAREN:
		p.advance()
		expression, err := p.parseExpr(preference(lexer.LPAREN))
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		p.advance()
		return expression, nil
	case lexer.LSQUARE:
		nodes, err := p.parseParameters(lexer.RSQUARE)
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		p.advance()
		return &ArrayNode{lexer.ARRAY_NODE, nodes}, nil
	case lexer.IDENTIFIER:
		identifier := p.token.Literal
		p.advance()
		if p.token.Type == lexer.LPAREN {
			params, err := p.parseParameters(lexer.RPAREN)
			if err != nil {
				return &ErrorNode{lexer.ERROR}, nil
			}
			p.advance()
			if p.token.Type == lexer.LPAREN {
				paramsB, err := p.parseParameters(lexer.RPAREN)
				if err != nil {
					return &ErrorNode{lexer.ERROR}, err
				}
				p.advance()
				return &FunctionCallNode{lexer.FUNCTION_CALL_NODE, identifier, params, paramsB}, nil
			}
			return &FunctionCallNode{lexer.FUNCTION_CALL_NODE, identifier, make([]Node, 0), params}, nil
		}
		return &IdentifierNode{lexer.IDENTIFIER, identifier}, nil
	case lexer.INT:
		value, err := strconv.Atoi(p.token.Literal)
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		p.advance()
		node := &IntNode{Type: lexer.INT_NODE, Value: value}
		if p.token.Type == lexer.IDENTIFIER {
			unit := p.token.Literal
			p.advance()
			return &UnitNode{lexer.UNIT_NODE, node, unit}, nil
		}
		return node, nil
	case lexer.FLOAT:
		value, err := strconv.ParseFloat(p.token.Literal, 64)
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		p.advance()
		node := &FloatNode{Type: lexer.INT_NODE, Value: value}
		if p.token.Type == lexer.IDENTIFIER {
			unit := p.token.Literal
			p.advance()
			return &UnitNode{lexer.UNIT_NODE, node, unit}, nil
		}
		return node, nil
	case lexer.STRING:
		value := p.token.Literal
		p.advance()
		return &StringNode{lexer.STRING_NODE, value}, nil
	}
	return &ErrorNode{lexer.ERROR}, errors.New("SyntaxError: parsePrefix() unsupported prefix:" + p.token.Literal)
}

func (p *Parser) parseInfix(left Node, operation string) (Node, error) {
	if !contains([]string{"EE", "NE", "LT", "GT", "LTE", "GTE", "ADD", "SUB", "MUL", "DIV", "MOD", "POW", "IN", "ARROW"}, p.token.Type) {
		return &ErrorNode{lexer.ERROR}, errors.New("SyntaxError: parseInfix() unsupported opperator:" + p.token.Literal)
	}
	p.advance()
	right, err := p.parseExpr(preference(operation) + 1) //-1
	if err != nil {
		return &ErrorNode{lexer.ERROR}, err
	}
	if operation == lexer.ARROW {
		operation = lexer.IN
	}
	return &BinOpNode{Type: lexer.BIN_OP_NODE, Left: left, Operation: operation, Right: right}, nil
}

func (p *Parser) parseParameters(terminate string) ([]Node, error) {
	parameters := make([]Node, 0)
	p.advance()
	if p.token.Type == lexer.RPAREN {
		return parameters, nil
	}
	for p.token.Type != terminate {
		if p.token.Type == lexer.EOF || p.token.Type == lexer.SEMICOLON {
			return make([]Node, 0), errors.New("SyntaxError: Unclosed parenthesis parseParameters()")
		}
		if p.token.Type != lexer.COMMA {
			param, err := p.parseExpr(0)
			if err != nil {
				return make([]Node, 0), err
			}
			parameters = append(parameters, param)
		} else {
			p.advance()
		}
	}
	return parameters, nil
}

func preference(tokenType string) int {
	var preferences = map[string]int{
		lexer.IN:     5,
		lexer.ARROW:  5,
		lexer.EE:     10,
		lexer.NE:     10,
		lexer.GT:     10,
		lexer.GTE:    10,
		lexer.LT:     10,
		lexer.LTE:    10,
		lexer.MOD:    15,
		lexer.ADD:    20,
		lexer.SUB:    20,
		lexer.MUL:    30,
		lexer.DIV:    30,
		lexer.POW:    40,
		lexer.LPAREN: 0,
		lexer.RPAREN: -1,
		lexer.COMMA:  -1,
	}

	if rbp, ok := preferences[tokenType]; ok {
		return rbp
	}
	return 0
}

func contains(array []string, element string) bool {
	for _, e := range array {
		if e == element {
			return true
		}
	}
	return false
}
