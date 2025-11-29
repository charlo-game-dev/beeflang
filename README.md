# Beeflang 🥩

**The first Beef Oriented Programming Language, in honor of the Church of Beef**

Beeflang is a **fully functional, Turing-complete** interpreted programming language built from scratch in Go. It features a Python-like syntax with beef-themed keywords, supporting functions, loops, conditionals, modules, and more. 

## Quick Start

**Prerequisites:** Go 1.21+

```bash
# Run a program
go run main.go examples/test.beef

# Run tests
go test ./...

# Dump tokens for debugging
go run main.go --dump-tokens examples/hello.beef
```

## Example Program

Here's a simple Beeflang program demonstrating the core features:

```beeflang
wrangle io

praise ChurchOfBeef():
  io.preach("What are you thankful for this Thanksgiving and why is it beef?")
  prep answer = io.input()

  if answer == "beef":
    io.preach("Braised be!")
  else:
    io.preach("You have been removed from Church of Beef")
  beef
beef
```

**See `examples/` directory for more complete programs** including Fibonacci, factorial, prime checking, and more!

## Language Reference

### Program Structure

Every Beeflang program must define a `ChurchOfBeef()` function as the entry point:

```beeflang
praise ChurchOfBeef():
  # Your code here
beef
```

The interpreter automatically calls `ChurchOfBeef()` when the program runs - you don't need to call it explicitly.

### Variables

```beeflang
prep x = 42              # Variable declaration
prep name = "Beef"       # Strings
prep is_tasty = true     # Booleans

x = 100                  # Reassignment (no 'prep')
```

### Data Types

- **Integers**: `42`, `-10`, `0`
- **Booleans**: `true`, `false`
- **Strings**: `"Hello, Beef!"` (double-quotes only)
- **Functions**: First-class values with closures

### Operators

**Arithmetic**: `+`, `-`, `*`, `/`, `%`
```beeflang
prep sum = 5 + 3      # 8
prep product = 4 * 2  # 8
prep modulo = 10 % 3  # 1
```

**Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=`
```beeflang
if x > 10:
  # do something
beef
```

**String Concatenation**: `+`
```beeflang
prep greeting = "Hello, " + "Beef!"
```

**Prefix Operators**: `-` (negation), `!` (logical not)
```beeflang
prep negative = -42
prep opposite = !true  # false
```

### Functions

```beeflang
praise add(x, y):
  serve x + y
beef

praise greet(name):
  io.preach("Hello, " + name)
  # No return - returns NULL
beef

praise ChurchOfBeef():
  prep result = add(5, 3)  # 8
  greet("Believer")
beef
```

Functions support:
- **Recursion**: Functions can call themselves
- **Closures**: Functions capture their surrounding environment
- **First-class**: Pass functions as values

### Conditionals

```beeflang
if condition:
  # consequence
else:
  # alternative
beef

# Without else
if x > 0:
  io.preach("Positive!")
beef
```

- Single `beef` closes the entire `if/else` block
- Conditions are "truthy" - `false` and `NULL` are falsy, everything else is truthy

### Loops

```beeflang
feast while condition:
  # loop body
beef

# Example: countdown
prep counter = 5
feast while counter > 0:
  io.preach(counter)
  counter = counter - 1
beef
```

### Modules

```beeflang
wrangle io  # Import the 'io' module

praise ChurchOfBeef():
  io.preach("Hello!")       # Print with newline
  prep name = io.input()    # Read line from stdin
beef
```

**Built-in modules:**
- `io.preach(value)` - Print to stdout with newline
- `io.input()` - Read line from stdin, returns string

### Comments

```beeflang
# Single-line comments only (Python-style)
prep x = 42  # Inline comments work too
```

### Keywords Reference

| Keyword | Purpose | Example |
|---------|---------|---------|
| `prep` | Variable declaration | `prep x = 42` |
| `praise` | Function declaration | `praise add(x, y):` |
| `serve` | Return from function | `serve x + y` |
| `if` / `else` | Conditionals | `if x > 0: ... else: ... beef` |
| `feast while` | While loop | `feast while x > 0: ... beef` |
| `beef` | Block terminator | Ends functions, loops, conditionals |
| `wrangle` | Import module | `wrangle io` |
| `true` / `false` | Boolean literals | `prep is_valid = true` |

### Syntax Rules

- **Indentation**: Recommended for readability (not enforced)
- **Newline-terminated**: Statements end at newlines (no semicolons)
- **Colons**: Required after function/loop/conditional headers
- **Block terminator**: Every block needs `beef` to close it

## More Examples

Check out `examples/` for complete programs:
- **`test.beef`** - Interactive I/O with conditionals
- **`fibonacci.beef`** - Iterative Fibonacci calculator
- **`factorial.beef`** - Recursive factorial
- **`prime_check.beef`** - Prime number checker with loops
- **`showcase.beef`** - GCD, summation, and power functions
- **`countdown.beef`** - Simple while loop demo

## Implementation Details

✅ **Fully Functional** - All core features implemented and tested!

The interpreter includes:
- **Lexer**: Character-by-character tokenization with line/column tracking
- **Parser**: Recursive descent with Pratt parsing for operator precedence
- **Evaluator**: Tree-walking interpreter with environment-based scoping
- **Type System**: Dynamic typing with runtime value objects
- **Module System**: Extensible module loader with dot notation

Built with **Test-Driven Development** (TDD) - comprehensive test coverage across all components.

### Architecture

```
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Output
```

See [CLAUDE.md](CLAUDE.md) for development workflow and [BEEFLANG_SPEC.md](BEEFLANG_SPEC.md) for the complete language specification.

## Why Beeflang?

This is a **learning project** to understand how interpreters work by building one from scratch. The beef theme makes it fun, but the implementation is serious - we follow interpreter design best practices and maintain clean, well-tested code.

Beeflang is **Turing-complete**, meaning it can compute anything computable. Try writing your own algorithms!

## License

This is a learning project. Do whatever you want with it. 🥩
