package lox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Lox struct {
	hadError bool
}

func NewLox() *Lox {
	return &Lox{
		hadError: false,
	}
}

func (l *Lox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file: ", err.Error())
		os.Exit(1)
	}

	l.run(string(bytes))
	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		l.run(strings.TrimSpace(line))
		l.hadError = false
	}
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source, l)
	tokens := scanner.scanTokens()

	// for _, token := range tokens {
	// 	fmt.Println(token.String())
	// }

	parser := NewParser(tokens, l)
	expr, err := parser.Parse()
	if err != nil {
		return
	}

	astPrinter := NewAstPrinter()
	fmt.Printf("AST: %s\n", astPrinter.Print(expr))
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("%v: [line %d] Error %s: %s\n", os.Stderr, line, where, message)
	l.hadError = true
}
