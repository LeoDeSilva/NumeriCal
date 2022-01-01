package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"numerical/evaluator"
	"numerical/lexer"
	"numerical/parser"
	"os"
	"strings"
)

func interpretProgram(program string, environment evaluator.Environment) error {
	l := lexer.NewLexer(program)
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
	fmt.Println(obj.String())
	return nil
}

func startRepl(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	periodicTable, _ := ioutil.ReadFile("/Users/ldesilva/Documents/Personal/Coding/Golang/NumeriCal/evaluator/periodicTable.json")

	var periodicTableJson map[string]interface{}
	json.Unmarshal([]byte(periodicTable), &periodicTableJson)

	environment := evaluator.Environment{Variables: make(map[string]evaluator.Object), Functions: make(map[string]*parser.FunctionDefenitionNode), PeriodicTable: periodicTableJson}

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
	if len(os.Args) > 1 {
		periodicTable, _ := ioutil.ReadFile("/Users/ldesilva/Documents/Personal/Coding/Golang/NumeriCal/evaluator/periodicTable.json")
		var periodicTableJson map[string]interface{}
		json.Unmarshal([]byte(periodicTable), &periodicTableJson)
		environment := evaluator.Environment{
			Variables:     make(map[string]evaluator.Object),
			Functions:     make(map[string]*parser.FunctionDefenitionNode),
			PeriodicTable: periodicTableJson,
		}

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
