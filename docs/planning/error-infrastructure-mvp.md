# Error Infrastructure - MVP Implementation Plan

## Goal
Build a minimal error handling foundation that supports clean error messages throughout the interpreter, making it easy to add detailed errors as we implement new features.

## Current State (Likely)
- Errors probably handled inconsistently:
  - Some: `return nil` (ambiguous)
  - Some: `fmt.Printf("error: ...")` (prints to stdout, hard to test)
  - Some: `panic()` (crashes interpreter)
  - Parser: might have error collection already

## MVP Target

### What We Want to Enable

```go
// In evaluator - clear, testable errors
if index >= len(array.Elements) {
    return newError("index out of bounds: tried to access index %d of array with length %d",
        index, len(array.Elements))
}

// In tests - easy to verify errors
result := testEval(input)
errObj, ok := result.(*object.Error)
assert.True(t, ok)
assert.Contains(t, errObj.Message, "index out of bounds")
```

### What User Sees (MVP)

```
Error at line 4, column 10 - index out of bounds: tried to access index 5 of array with length 3
```

**Future:** Add source snippet with caret pointing to error location.

---

## Implementation Layers

### 1. Object Layer - Error Type

**Add to `internal/object/object.go`:**

```go
// Error represents a runtime error in Beeflang.
// When evaluation encounters an error, it returns an Error object
// instead of panicking or returning nil.
//
// Location information (Line, Column, File) is included from the start
// because Token already tracks this data and it's much easier to thread
// through during implementation than to retrofit later.
type Error struct {
    Message string
    Line    int    // Line number where error occurred (from Token)
    Column  int    // Column number where error occurred (from Token)
    File    string // Source file path (empty string if not from file)
}

func (e *Error) Type() string {
    return "ERROR"
}

func (e *Error) Inspect() string {
    if e.File != "" {
        return fmt.Sprintf("Error at %s:%d:%d - %s",
            e.File, e.Line, e.Column, e.Message)
    }
    if e.Line > 0 {
        return fmt.Sprintf("Error at line %d, column %d - %s",
            e.Line, e.Column, e.Message)
    }
    return "Error: " + e.Message
}
```

**Helper function in `internal/evaluator/evaluator.go`:**

```go
// newError creates an Error object with a formatted message and location.
// Usage: return newError(node.Token, "unexpected type: got %s", obj.Type())
func newError(tok token.Token, format string, a ...interface{}) *object.Error {
    return &object.Error{
        Message: fmt.Sprintf(format, a...),
        Line:    tok.Line,
        Column:  tok.Column,
        // File is set by main.go when running from a file
    }
}
```

**Tests:**
```go
func TestErrorObjectType(t *testing.T)
func TestErrorObjectInspect(t *testing.T)
```

---

### 2. Evaluator Layer - Error Handling Pattern

**Pattern to follow throughout evaluator:**

```go
// Before (bad):
if someCondition {
    return nil  // ambiguous!
}

// After (good):
if someCondition {
    return newError("clear description of what went wrong")
}
```

**Error checking pattern:**

```go
// When evaluating sub-expressions, check for errors
val := Eval(node.SomeExpression, env)
if isError(val) {
    return val  // propagate error up the call stack
}

// Example with location:
if left.Type() != "INTEGER" {
    return newError(node.Token, "type mismatch: %s + %s",
        left.Type(), right.Type())
}

// Helper function:
func isError(obj object.Object) bool {
    return obj != nil && obj.Type() == "ERROR"
}
```

**Files to modify:**
- `internal/evaluator/evaluator.go` - add `newError()` and `isError()` helpers

**Tests:**
```go
func TestErrorPropagation(t *testing.T)  // Errors bubble up correctly
```

---

### 3. Parser Layer - Error Collection

**Check current state first:**
- Does parser already collect errors?
- Look for `p.errors []string` or similar

**If parser has error collection (likely):**
- ✅ Already done, just use it consistently
- Make sure errors are descriptive

**If parser doesn't have error collection:**

Add to parser struct:
```go
type Parser struct {
    // ... existing fields ...
    errors []string
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) addError(msg string) {
    p.errors = append(p.errors, msg)
}
```

**Tests:**
- Parser tests should check `len(p.Errors())` for syntax errors

---

### 4. Main/REPL Layer - Error Display

**In `main.go` (or wherever you run programs):**

```go
// After evaluation:
result := evaluator.Eval(program, env)

// Check if result is an error
if result != nil && result.Type() == "ERROR" {
    fmt.Fprintf(os.Stderr, "%s\n", result.Inspect())
    os.Exit(1)
}

// Otherwise print result
if result != nil {
    fmt.Println(result.Inspect())
}
```

**For parser errors:**

```go
parser := parser.New(lexer)
program := parser.ParseProgram()

if len(parser.Errors()) > 0 {
    for _, msg := range parser.Errors() {
        fmt.Fprintf(os.Stderr, "Parse error: %s\n", msg)
    }
    os.Exit(1)
}
```

---

## Common Error Scenarios to Handle

### Runtime Errors (Evaluator)

1. **Type errors:**
   ```go
   return newError(node.Token, "type mismatch: %s + %s", left.Type(), right.Type())
   ```

2. **Undefined variables:**
   ```go
   return newError(node.Token, "identifier not found: %s", node.Value)
   ```

3. **Invalid operations:**
   ```go
   return newError(node.Token, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
   ```

4. **Index errors (future - arrays):**
   ```go
   return newError(node.Token, "index out of bounds: index %d, length %d", idx, len)
   return newError(node.Token, "index operator not supported: %s", obj.Type())
   ```

5. **Wrong number of arguments:**
   ```go
   return newError(node.Token, "wrong number of arguments: expected %d, got %d",
       expectedLen, gotLen)
   ```

### Parse Errors (Parser)

1. **Unexpected tokens:**
   ```
   expected next token to be ')', got 'IDENT' instead
   ```

2. **Invalid syntax:**
   ```
   no prefix parse function for token type 'RBRACE'
   ```

---

## Testing Strategy

### Unit Tests (Evaluator)

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
        {"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"}, // stops at first error
        {"-true", "unknown operator: -BOOLEAN"},
        {"foobar", "identifier not found: foobar"},
    }

    for _, tt := range tests {
        evaluated := testEval(tt.input)

        errObj, ok := evaluated.(*object.Error)
        assert.True(t, ok, "expected error object")
        assert.Contains(t, errObj.Message, tt.expected)
    }
}
```

### Integration Tests

Run `.beef` files and verify error output:
```bash
$ go run main.go examples/error_test.beef
Error: identifier not found: undefinedVar
```

---

## Implementation Order (TDD)

### Phase 1: Object + Helpers (30 min)
1. Add `object.Error` type
2. Add `newError()` helper to evaluator
3. Add `isError()` helper to evaluator
4. Write tests for Error object

**Milestone:** Can create and inspect error objects

### Phase 2: Evaluator Integration (45 min)
1. Find all places that return `nil` or panic on errors
2. Replace with `return newError("...")`
3. Add error checks after evaluating sub-expressions
4. Write comprehensive error handling tests

**Milestone:** Evaluator returns errors instead of nil/panic

### Phase 3: Parser Errors (15 min)
1. Verify parser has error collection (probably already does)
2. Make error messages descriptive if needed
3. Add tests for common syntax errors

**Milestone:** Parser errors are consistent and clear

### Phase 4: Main Integration (15 min)
1. Update `main.go` to check for errors
2. Print errors to stderr
3. Exit with non-zero status on error

**Milestone:** User sees clear error messages when running programs

---

## Error Checklist for Future Features

When adding new features (like arrays), follow this pattern:

```go
// 1. Evaluate sub-expressions and check for errors
left := Eval(node.Left, env)
if isError(left) {
    return left
}

// 2. Validate types
if left.Type() != "ARRAY" {
    return newError("index operator not supported: %s", left.Type())
}

// 3. Validate operations
if index < 0 || index >= len(array.Elements) {
    return newError("index out of bounds: index %d, length %d",
        index, len(array.Elements))
}

// 4. Return result or error
return result
```

---

## Success Criteria (MVP)

✅ `object.Error` type exists with Message, Line, Column, File fields
✅ Error implements `Object` interface
✅ `Inspect()` displays location information when available
✅ `newError()` helper takes token and creates error with location
✅ `isError()` helper for checking if object is error
✅ Evaluator returns errors instead of nil/panic
✅ Errors propagate up the call stack (don't get swallowed)
✅ Parser has descriptive error messages
✅ Main displays errors to stderr and exits with status 1
✅ Can write tests that verify error messages and locations
✅ Common errors have helpful messages with line/column info

---

## Future Enhancements (Post-MVP)

### Short Term (After Arrays)

1. **Source Snippets**
   - Read source file and display the offending line
   - Add caret pointing to error location:
     ```
     Error at line 4, column 10 - index out of bounds
       4 |   prep x = arr[5]
                          ^
     ```

2. **Stack Traces**
   - Track function call stack
   - Show stack trace on error:
     ```
     Error: index out of bounds
       at fibonacci (line 10)
       at main (line 5)
     ```

3. **Source Snippets**
   - Show the offending line with caret:
     ```
     4 |   prep x = arr[5]
                        ^
     ```

### Medium Term

4. **Error Recovery**
   - Parser: Continue parsing after errors (collect multiple errors)
   - Helps users fix multiple syntax errors at once

5. **Warnings**
   - Separate `Warning` type for non-fatal issues
   - Unused variables, unreachable code, etc.

6. **Error Codes**
   - Assign codes to error types (E001, E002, etc.)
   - Link to documentation for each error

### Long Term

7. **Rich Error Messages**
   - Suggestions: "Did you mean 'feast while' instead of 'while'?"
   - Type hints: "Expected INTEGER, got BOOLEAN. Consider using `if` instead of `+`"

8. **Language Server Protocol (LSP)**
   - Real-time error checking in editor
   - Squiggly underlines for errors
   - Part of VSCode extension work

---

## Examples

### Before MVP (Inconsistent)
```go
// Various places in code:
if something {
    return nil  // What went wrong?
}

if otherThing {
    fmt.Println("Error: bad thing")  // Prints to stdout, can't test
    return nil
}

if yetAnother {
    panic("oh no")  // Crashes interpreter
}
```

### After MVP (Clean)
```go
// Consistent error pattern:
if !isValidType(obj) {
    return newError("type mismatch: expected %s, got %s",
        expected, obj.Type())
}

// In tests:
result := testEval(input)
assert.IsType(t, &object.Error{}, result)
assert.Contains(t, result.(*object.Error).Message, "type mismatch")
```

### After Source Snippets (Future)
```go
// MVP already has location in errors
// Future: Display source context by reading file

// User sees:
// Error at line 4, column 10 - index out of bounds
//   4 |   prep x = arr[5]
//                      ^
```

---

## Notes

- **Keep it simple:** MVP is just error objects + consistent handling
- **No over-engineering:** Don't build a fancy error framework yet
- **Test-friendly:** Main goal is making errors easy to test
- **Foundation for growth:** Easy to add location info later
- **Learning opportunity:** See how real interpreters handle errors

## Estimated Time
- MVP implementation: **1.5-2 hours**
- Testing: **30 minutes**
- Total: **~2-2.5 hours**

Then you're ready to implement arrays with clean error handling from day one!
