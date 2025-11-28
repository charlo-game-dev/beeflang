# CLAUDE.md - Beeflang Development Guide

## Project Overview
Beeflang is a meme language interpreter written in Go, honoring the Church of Beef.

**This is a learning project** - the goal is to understand how interpreters work by building one from scratch. We move slowly, understand each piece through testing, and let the implementation emerge from TDD rather than copying patterns.

## Architecture

### Interpreter Pipeline
```
Source Code (.beef) → Lexer → Tokens → Parser → AST → Evaluator → Output
```

### Components

1. **Lexer** (`internal/lexer/`)
   - Reads source code character by character
   - Produces stream of tokens
   - Example: `cut x = 42` → `[CUT, IDENT("x"), ASSIGN, INT(42)]`

2. **Token** (`internal/token/`)
   - Defines all token types (keywords, operators, literals)
   - Minimal data structure representing lexical elements

3. **Parser** (`internal/parser/`)
   - Consumes tokens from lexer
   - Builds Abstract Syntax Tree (AST)
   - Validates syntax rules

4. **AST** (`internal/ast/`)
   - Tree representation of program structure
   - Nodes for expressions, statements, literals, etc.

5. **Evaluator** (`internal/evaluator/`)
   - Walks the AST
   - Executes the program
   - Manages runtime state (variables, scope, etc.)

6. **Object/Value System** (`internal/object/`)
   - Runtime value representations (int, bool, functions)
   - Type system implementation

7. **Standard Library** (`internal/stdlib/`)
   - Built-in functions like `preach()`
   - Kept separate from core language

## Development Workflow

### TDD Cycle (Strict)

**Phase 1: Scaffold Empty Tests**
- Write test function names with `assert true == false` placeholder
- Review structure and test names with Claude
- Iterate until test coverage plan is agreed upon

**Phase 2: Write Real Failing Tests**
- Replace scaffolds with actual test logic
- Tests should fail because implementation doesn't exist yet (RED)
- Review test assertions and edge cases

**Phase 3: Implementation**
- Write minimal code to make tests pass (GREEN)
- Run tests frequently

**Phase 4: Refactor**
- Clean up code while keeping tests green
- Improve naming, structure, performance

### Running the Project

**Early Stage (Now):**
```bash
# Run interpreter on a .beef file
go run main.go examples/hello.beef

# Run tests
go test ./...

# Run specific package tests
go test ./internal/lexer -v
```

**Future (CLI Tool):**
```bash
beeflang run hello.beef
```

### Example Test Flow

```go
// Phase 1: Scaffold
func TestLexerTokenizesVariableDeclaration(t *testing.T) {
    assert.True(t, false, "TODO: implement")
}

// Phase 2: Real failing test
func TestLexerTokenizesVariableDeclaration(t *testing.T) {
    input := "cut x = 42"
    lexer := New(input)

    expectedTokens := []token.Token{
        {Type: token.CUT, Literal: "cut"},
        {Type: token.IDENT, Literal: "x"},
        {Type: token.ASSIGN, Literal: "="},
        {Type: token.INT, Literal: "42"},
    }

    for i, expected := range expectedTokens {
        tok := lexer.NextToken()
        assert.Equal(t, expected.Type, tok.Type)
        assert.Equal(t, expected.Literal, tok.Literal)
    }
}

// Phase 3: Implement lexer.NextToken() until test passes
// Phase 4: Refactor if needed
```

## Implementation Order

1. **Token + Lexer** - Break source into tokens
2. **AST Nodes** - Define tree structure
3. **Parser** - Build AST from tokens
4. **Object System** - Runtime values
5. **Evaluator** - Execute AST
6. **Standard Library** - Built-in functions like `preach()`
7. **CLI** - Polish the interface

## Language Spec Quick Reference

### Keywords
- `praise` - function declaration
- `beef` - block terminator
- `feast while` - while loop
- `if` / `else` - conditionals
- `cut` - variable declaration
- `serve` - return statement
- `genesis` - entry point (main function)
- `preach` - print/output (stdlib function)

### Types
- `int` - integers
- `bool` - true/false
- Dynamically typed (inferred)

### Syntax Rules
- Newline-terminated statements (no semicolons)
- Block-level scoping
- Colons after function/loop/conditional headers
- Comments with `#`

### Example Program
```beeflang
praise genesis():
   cut x = 42
   preach(x)

   if x > 0:
      cut y = 5
   beef
beef
```

## Error Handling Philosophy

Keep it simple for now:
- Basic error messages: "syntax error at line X"
- Don't over-engineer error recovery
- Focus on happy path first
- Can enhance later

## Testing Philosophy

- Focus on **happy path** unit tests
- Test core functionality at each layer
- Don't test every edge case initially
- Build confidence in the main workflow

## File Organization

```
beeflang/
├── main.go                 # Entry point
├── internal/
│   ├── token/             # Token definitions
│   ├── lexer/             # Lexical analysis
│   ├── ast/               # AST node definitions
│   ├── parser/            # Syntax analysis
│   ├── object/            # Runtime value system
│   ├── evaluator/         # Execution engine
│   └── stdlib/            # Standard library (preach, etc.)
├── examples/              # Sample .beef programs
└── BEEFLANG_SPEC.md       # Language specification
```

## Development Notes

- This is a **learning project** - keep it simple
- Aim for **minimal Turing completeness** first
- Room to grow with features later
- Follow Go idioms and conventions
- Keep packages small and focused
