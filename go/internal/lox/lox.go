package lox

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
	interpreter     *Interpreter
}

func NewLox() *Lox {
	return &Lox{
		hadError:        false,
		hadRuntimeError: false,
		interpreter:     NewInterpreter(),
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
	if l.hadRuntimeError {
		os.Exit(70)
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
		l.hadRuntimeError = false
	}
}

func (l *Lox) run(source string) {
	// Scanning
	scanner := NewScanner(source)
	tokens, scanErrors := scanner.ScanTokens()

	// Report scan errors
	for _, err := range scanErrors {
		l.reportError(err)
	}
	if len(scanErrors) > 0 {
		l.hadError = true
		return
	}

	// Parsing
	parser := NewParser(tokens)
	expression, parseErrors := parser.Parse()

	// Report parse errors
	for _, err := range parseErrors {
		l.reportError(err)
	}
	if len(parseErrors) > 0 {
		l.hadError = true
		return
	}

	// Interpreting
	result, err := l.interpreter.Interpret(expression)
	if err != nil {
		l.reportRuntimeError(err)
		l.hadRuntimeError = true
		return
	}

	fmt.Println(result)
}

func (l *Lox) reportError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}

func (l *Lox) reportRuntimeError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}
