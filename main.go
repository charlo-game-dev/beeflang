package main

import (
	"fmt"
	"os"

	"github.com/elitwilson/beeflang/internal/evaluator"
	"github.com/elitwilson/beeflang/internal/lexer"
	"github.com/elitwilson/beeflang/internal/object"
	"github.com/elitwilson/beeflang/internal/parser"
	"github.com/elitwilson/beeflang/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go <file.beef>")
		fmt.Println("  go run main.go --dump-tokens <file.beef>")
		os.Exit(1)
	}

	// Check for --dump-tokens flag
	dumpTokens := false
	filename := os.Args[1]

	if os.Args[1] == "--dump-tokens" {
		if len(os.Args) < 3 {
			fmt.Println("Error: --dump-tokens requires a filename")
			os.Exit(1)
		}
		dumpTokens = true
		filename = os.Args[2]
	}

	// Read source file
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Dump tokens mode
	if dumpTokens {
		l := lexer.New(string(source))
		fmt.Printf("Tokens for %s:\n", filename)
		fmt.Println("---")
		for {
			tok := l.NextToken()
			fmt.Printf("%-15s %-10s (line %d, col %d)\n", tok.Type, tok.Literal, tok.Line, tok.Column)
			if tok.Type == token.EOF {
				break
			}
		}
		return
	}

	// Normal interpreter mode - run the program!
	l := lexer.New(string(source))
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parser errors
	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Printf("  %s\n", msg)
		}
		os.Exit(1)
	}

	// Evaluate the program (this loads all function/variable declarations)
	env := object.NewEnvironment()
	evaluator.Eval(program, env)

	// Auto-call ChurchOfBeef() if it exists (entry point function)
	if entryPoint, ok := env.Get("ChurchOfBeef"); ok {
		if fn, ok := entryPoint.(*object.Function); ok {
			// Create new environment for ChurchOfBeef() execution
			entryEnv := object.NewEnclosedEnvironment(fn.Env)
			// Execute ChurchOfBeef() body
			evaluator.Eval(fn.Body, entryEnv)
		} else {
			fmt.Println("Error: ChurchOfBeef is not a function")
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: no ChurchOfBeef() entry point function found")
		os.Exit(1)
	}
}
