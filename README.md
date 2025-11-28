# Beeflang 🥩

A Turing-complete programming language honoring the Church of Beef

## What is this?

Beeflang is a **fully functional, Turing-complete** interpreted programming language built from scratch in Go as a learning project. It features beef-themed keywords, Python-like syntax, and supports functions, loops, conditionals, and recursion.

## Example Programs

### Hello World
```beeflang
# Simple program - returns 42
prep answer = 42
answer
```

### Fibonacci (Iterative)
```beeflang
praise fibonacci(n):
   if n <= 1:
      serve n
   beef

   prep a = 0
   prep b = 1
   prep i = 2

   feast while i <= n:
      prep temp = a + b
      a = b
      b = temp
      i = i + 1
   beef

   serve b
beef

fibonacci(10)  # Returns 55
```

### Factorial (Recursive)
```beeflang
praise factorial(n):
   if n <= 1:
      serve 1
   beef
   serve n * factorial(n - 1)
beef

factorial(5)  # Returns 120
```

## Quick Start

**Prerequisites:** Go 1.21+

```bash
# Run a program
go run main.go examples/fibonacci.beef

# Run tests
go test ./...

# Dump tokens for debugging
go run main.go --dump-tokens examples/hello.beef
```

## Language Features

✅ **Turing Complete** - Can compute anything computable!

- **Variables**: `prep x = 42` (declaration), `x = 10` (reassignment)
- **Functions**: First-class with closures and recursion
- **Conditionals**: `if`/`else` statements
- **Loops**: `feast while` for iteration
- **Types**: Integers, booleans, strings, functions
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), comparison, logic

## Keywords

- `prep` - variable declaration
- `praise` - function declaration
- `serve` - return statement
- `if` / `else` - conditionals
- `feast while` - while loop
- `beef` - block terminator
- `true` / `false` - boolean literals

## Status

✅ **Fully Functional** - All core features implemented and tested!

The interpreter includes:
- Complete lexer with position tracking
- Recursive descent parser with Pratt parsing for expressions
- Tree-walking evaluator with environment-based scoping
- Comprehensive test suite with TDD workflow

See [CLAUDE.md](CLAUDE.md) for development workflow and [BEEFLANG_SPEC.md](BEEFLANG_SPEC.md) for full language specification.

## Examples

Check out `examples/` for more programs:
- `fibonacci.beef` - Iterative Fibonacci calculator
- `factorial.beef` - Recursive factorial
- `prime_check.beef` - Prime number checker
- `showcase.beef` - Comprehensive feature demo

## License

This is a learning project. Do whatever you want with it.
