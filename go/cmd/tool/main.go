package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shresth72/lox/internal/tool"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tool <command> [args]")
		os.Exit(64)
	}

	command := os.Args[1]

	switch command {
	case "generate_ast":
		runGenerateAst()
	case "print_ast":
		runPrintAst()
	default:
		fmt.Printf("Tool: %s not supported\n", command)
	}
}

func runGenerateAst() {
	var outputDir string
	if len(os.Args) != 3 {
		projectRoot, err := findProjectRoot()
		if err != nil {
			fmt.Println("Error: could not determine project root:", err)
			os.Exit(1)
		}
		outputDir = filepath.Join(projectRoot, "internal", "lox")
	} else {
		outputDir = os.Args[2]
	}

	ast := tool.NewAST()
	if err := ast.GenerateAST(outputDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func runPrintAst() {
	fmt.Printf("Example AST: %s\n", tool.ExampleAst())
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find project root (go.mod not found)")
}
