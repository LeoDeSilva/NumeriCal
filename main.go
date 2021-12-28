package main

import (
	"bufio"
	"fmt"
	"io"
	"numerical/lexer"
	"numerical/parser"
	"os"
)

func startRepl(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

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

		l := lexer.NewLexer(line)
		tokens, err := l.Lex()
		if err != nil {
			return err
		}
		fmt.Println(tokens)

		p := parser.NewParser(tokens)
		ast, err := p.Parse()
		if err != nil {
			return err
		}
		fmt.Println(ast)
	}

	return nil
}

func run() error {
	err := startRepl(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
