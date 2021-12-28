package parser

import (
	"errors"
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

type IdentifierNode struct {
	Type       string
	Identifier string
}

func (i *IdentifierNode) Eval() Node     { return i }
func (i *IdentifierNode) String() string { return "{IDENTIFIER_NODE:" }

type IntNode struct {
	Type  string
	Value int
}

func (i *IntNode) Eval() Node     { return i }
func (i *IntNode) String() string { return "{INT_NODE:" + strconv.Itoa(i.Value) + "}" }

type BinOpNode struct {
	Type      string
	Left      Node
	Operation string
	Right     Node
}

func (b *BinOpNode) Eval() Node { return b }
func (b *BinOpNode) String() string {
	return "{BIN_OP_NODE:" + b.Left.String() + b.Operation + b.Right.String() + "}"
}

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
		peekRbp = preference(p.peekToken().Type)
	}
	return left, nil
}

func (p *Parser) parsePrefix() (Node, error) {
	//Unit -> INT/FLOAT IDENTIFIER
	//PAREN -> ( Expr )
	//FUNCTION -> ID ( ... )
	switch p.token.Type {
	case lexer.INT:
		value, err := strconv.Atoi(p.token.Literal)
		if err != nil {
			return &ErrorNode{lexer.ERROR}, err
		}
		p.advance()
		return &IntNode{Type: lexer.INT_NODE, Value: value}, nil
	case lexer.LPAREN:

	}
	return &ErrorNode{lexer.ERROR}, errors.New("SyntaxError: parsePrefix() unsupported prefix:" + p.token.Literal)
}

func (p *Parser) parseInfix(left Node, operation string) (Node, error) {
	if !contains([]string{"EE", "NE", "LT", "GT", "LTE", "GTE", "ADD", "SUB", "MUL", "DIV", "MOD", "POW"}, p.token.Type) {
		return &ErrorNode{lexer.ERROR}, errors.New("SyntaxError: parseInfix() unsupported opperator:" + p.token.Literal)
	}
	p.advance()
	right, err := p.parseExpr(preference(operation) + 1) //-1
	if err != nil {
		return &ErrorNode{lexer.ERROR}, err
	}

	return &BinOpNode{Type: lexer.BIN_OP_NODE, Left: left, Operation: operation, Right: right}, nil
}

func preference(tokenType string) int {
	var preferences = map[string]int{
		lexer.EE:     10,
		lexer.NE:     10,
		lexer.GT:     10,
		lexer.GTE:    10,
		lexer.LT:     10,
		lexer.LTE:    10,
		lexer.ADD:    20,
		lexer.SUB:    20,
		lexer.MUL:    30,
		lexer.DIV:    30,
		lexer.POW:    40,
		lexer.IN:     50,
		lexer.ARROW:  50,
		lexer.LPAREN: 0,
		lexer.EOF:    0,
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
