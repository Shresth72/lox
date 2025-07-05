package tool

import (
	"fmt"
	"os"
	"strings"
)

type AST struct{}

func NewAST() *AST {
	return &AST{}
}

func (ast *AST) GenerateAST(outputDir string) error {
	err := ast.defineAst(outputDir, "Expr", []string{
		"Binary: Expr left, Token operator, Expr right",
		"Grouping: Expr expression",
		"Literal: any value",
		"Unary: Token operator, Expr right",
	})
	return err
}

func (ast *AST) defineAst(outputDir, baseName string, types []string) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	path := outputDir + "/" + baseName + ".go"

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	write := func(s string) {
		file.WriteString(s + "\n")
	}

	write("package lox\n")

	write(
		fmt.Sprintf(
			"type %s interface {\n\tAccept(v %sVisitor) interface{}\n}",
			baseName,
			baseName,
		),
	)

	ast.defineVisitor(file, baseName, types)
	ast.defineTypes(file, baseName, types)

	return nil
}

func (ast *AST) defineTypes(file *os.File, baseName string, types []string) {
	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		fields := strings.TrimSpace(parts[1])

		ast.defineType(file, baseName, className, fields)
	}
}

func (ast *AST) defineType(file *os.File, baseName, className, fieldList string) {
	write := func(s string) {
		file.WriteString(s + "\n")
	}

	// Struct definition
	write(fmt.Sprintf("\ntype %s struct {", className))
	fields := strings.Split(fieldList, ",")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		fieldParts := strings.SplitN(field, " ", 2)
		fieldType := strings.TrimSpace(fieldParts[0])
		fieldName := strings.TrimSpace(fieldParts[1])
		write(fmt.Sprintf("    %s %s", capitalize(fieldName), fieldType))
	}
	write("}")

	// Constructor
	write(
		fmt.Sprintf(
			"\nfunc New%s(%s) *%s {",
			className,
			formatConstructorParams(fields),
			className,
		),
	)
	write("    return &" + className + "{")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		fieldName := strings.TrimSpace(strings.SplitN(field, " ", 2)[1])
		write(fmt.Sprintf("        %s: %s,", capitalize(fieldName), fieldName))
	}
	write("    }")
	write("}")

	// Accept method
	write(fmt.Sprintf("\nfunc (e *%s) Accept(v %sVisitor) interface{} {", className, baseName))
	write(fmt.Sprintf("    return v.Visit%s%s(e)", className, baseName))
	write("}")
}

func (ast *AST) defineVisitor(file *os.File, baseName string, types []string) {
	write := func(s string) {
		file.WriteString(s + "\n")
	}

	write(fmt.Sprintf("\ntype %sVisitor interface {", baseName))
	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		write(fmt.Sprintf("    Visit%s%s(expr *%s) interface{}", className, baseName, className))
	}
	write("}")
}

func capitalize(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func formatConstructorParams(fields []string) string {
	params := make([]string, len(fields))
	for i, field := range fields {
		field = strings.TrimSpace(field)
		parts := strings.SplitN(field, " ", 2)
		fieldType := strings.TrimSpace(parts[0])
		fieldName := strings.TrimSpace(parts[1])
		params[i] = fmt.Sprintf("%s %s", fieldName, fieldType)
	}
	return strings.Join(params, ", ")
}
