package main

import (
	"bufio"
	"fmt"
	"io"
	"numerical/evaluator"
	"numerical/lexer"
	"numerical/parser"
	"os"
	"strings"
)

func interpretProgram(program string, environment evaluator.Environment) error {
	l := lexer.NewLexer(strings.TrimSpace(program))
	tokens, err := l.Lex()
	if err != nil {
		return err
	}

	p := parser.NewParser(tokens)
	ast, err := p.Parse()
	if err != nil {
		return err
	}

	obj, err := evaluator.Eval(&ast, environment)
	if err != nil {
		return err
	}

	objString := obj.String()
	if objString != "" {
		fmt.Println(objString)
	}
	return nil
}

func startRepl(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	environment := evaluator.GenerateEnvironment()

	for {
		fmt.Fprintf(out, ">>")
		scanned := scanner.Scan()

		if !scanned {
			continue
		}

		line := scanner.Text()

		if line == "quit" {
			break
		}

		err := interpretProgram(line, environment)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

	return nil
}

func main() {
	evaluator.DefineUnits()
	if len(os.Args) > 1 {
		environment := evaluator.GenerateEnvironment()
		program := strings.Join(os.Args[1:], " ")
		err := interpretProgram(program, environment)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		startRepl(os.Stdin, os.Stdout)
	}
}
