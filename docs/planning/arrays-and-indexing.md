# Arrays and Indexing - Implementation Plan

## Goal
Add array literals, array operations, and index expressions to Beeflang.

## Target Syntax

```beeflang
wrangle io

praise ChurchOfBeef():
   # Array literals
   prep numbers = [1, 2, 3, 4, 5]
   prep mixed = [42, true, "beef"]
   prep empty = []

   # Index access (reading)
   prep first = numbers[0]
   prep last = numbers[4]

   # Index assignment (writing)
   numbers[0] = 99

   # String indexing (bonus - reuses same infrastructure)
   prep greeting = "hello"
   prep first_char = greeting[0]  # "h"

   io.preach(first)
   io.preach(numbers[0])
beef
```

## Implementation Layers

### 1. Token Layer (`internal/token/`)
**New tokens needed:**
- `LBRACKET` = `[`
- `RBRACKET` = `]`
- `COMMA` = `,` (if not already present)

**Files to modify:**
- `internal/token/token.go`

**Tests:**
- `internal/token/token_test.go` (if exists, otherwise lexer tests cover this)

---

### 2. Lexer Layer (`internal/lexer/`)
**Changes needed:**
- Recognize `[`, `]`, `,` as single-character tokens
- Already handles strings, so no changes needed for string indexing

**Files to modify:**
- `internal/lexer/lexer.go` - add cases in `NextToken()` switch

**Tests to write:**
```go
func TestLexerTokenizesArrayLiterals(t *testing.T)
func TestLexerTokenizesIndexExpressions(t *testing.T)
```

---

### 3. AST Layer (`internal/ast/`)
**New node types needed:**

```go
// ArrayLiteral represents an array literal: [1, 2, 3]
type ArrayLiteral struct {
    Token    token.Token  // The '[' token
    Elements []Expression // List of expressions inside brackets
}

// IndexExpression represents array/string indexing: arr[0]
type IndexExpression struct {
    Token token.Token  // The '[' token
    Left  Expression   // The array/string being indexed
    Index Expression   // The index expression (could be any expression)
}
```

**Files to modify:**
- `internal/ast/ast.go` - add new node types
- Both implement `Expression` interface

**Tests:**
- `internal/ast/ast_test.go` - test `String()` methods

---

### 4. Parser Layer (`internal/parser/`)
**New parsing functions needed:**

```go
func (p *Parser) parseArrayLiteral() Expression
func (p *Parser) parseIndexExpression(left Expression) Expression
```

**Changes to existing code:**
- Register `token.LBRACKET` as prefix parse function (for array literals)
- Register `token.LBRACKET` as infix parse function (for index expressions)
- Handle comma-separated lists in array literals

**Precedence:**
- Index expressions should be high precedence (same as call expressions)

**Files to modify:**
- `internal/parser/parser.go`

**Tests to write:**
```go
func TestParsingArrayLiterals(t *testing.T)
func TestParsingEmptyArrayLiterals(t *testing.T)
func TestParsingArrayLiteralsWithExpressions(t *testing.T)
func TestParsingIndexExpressions(t *testing.T)
func TestParsingNestedIndexExpressions(t *testing.T) // arr[0][1]
```

---

### 5. Object Layer (`internal/object/`)
**New object type needed:**

```go
// Array represents an array at runtime
type Array struct {
    Elements []Object  // Can hold any mix of types
}

func (a *Array) Type() string {
    return "ARRAY"
}

func (a *Array) Inspect() string {
    // Return string like "[1, 2, 3]"
}
```

**Files to modify:**
- `internal/object/object.go`

**Tests:**
- `internal/object/object_test.go` - test `Type()` and `Inspect()`

---

### 6. Evaluator Layer (`internal/evaluator/`)
**New evaluation functions needed:**

```go
func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object
func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object
```

**Logic for `evalArrayLiteral`:**
1. Evaluate each element expression
2. Collect results into `[]Object`
3. Return `&object.Array{Elements: elements}`

**Logic for `evalIndexExpression`:**
1. Evaluate the left side (should be Array or String)
2. Evaluate the index (should be Integer)
3. Check bounds
4. Return element at index
5. Handle both `object.Array` and `object.String`

**Index assignment (bonus):**
- Will need to handle `arr[0] = value` in assignment statements
- This is trickier - might defer to phase 2

**Files to modify:**
- `internal/evaluator/evaluator.go`

**Tests to write:**
```go
func TestEvalArrayLiterals(t *testing.T)
func TestEvalArrayIndexExpressions(t *testing.T)
func TestEvalArrayIndexOutOfBounds(t *testing.T)
func TestEvalStringIndexExpressions(t *testing.T)
func TestEvalStringIndexOutOfBounds(t *testing.T)
func TestArraysWithMixedTypes(t *testing.T)
```

---

### 7. Builtins (Nice to Have)
**Useful array builtins to add:**

```go
len(arr)       // Return length of array or string
push(arr, val) // Append to array (returns new array? or mutates?)
pop(arr)       // Remove last element
```

**Files to modify:**
- `internal/evaluator/builtins.go` (or wherever builtins live)
- Add to builtin registry

**Philosophy question:** Immutable or mutable arrays?
- **Immutable:** `push(arr, 5)` returns new array
- **Mutable:** `push(arr, 5)` modifies array in place

Recommendation: Start immutable (simpler), add mutation later if needed.

---

## Implementation Order (TDD)

### Phase 1: Array Literals (Read-only)
1. Token + Lexer: Add `[`, `]`, `,` tokens
2. AST: Add `ArrayLiteral` node
3. Parser: Parse `[1, 2, 3]`
4. Object: Add `Array` type
5. Evaluator: Evaluate array literals

**Milestone:** Can create and print arrays

### Phase 2: Index Reading
1. AST: Add `IndexExpression` node
2. Parser: Parse `arr[0]`
3. Evaluator: Evaluate index expressions (arrays only)

**Milestone:** Can read from arrays

### Phase 3: String Indexing (Bonus)
1. Evaluator: Extend `evalIndexExpression` to handle `object.String`

**Milestone:** Can index strings like `"hello"[0]`

### Phase 4: Builtins
1. Add `len()` builtin
2. Add `push()` builtin (immutable)

**Milestone:** Can get length and append to arrays

### Phase 5: Index Assignment (Optional - More Complex)
1. Parser: Handle `arr[0] = value` in assignment
2. Evaluator: Implement index assignment

**Milestone:** Can modify array elements

---

## Edge Cases to Consider

1. **Out of bounds access:** Return error or null?
2. **Negative indices:** Error or Python-style wraparound?
3. **Index non-integer:** Clear error message
4. **Index non-indexable:** Clear error message
5. **Empty arrays:** `[]` should work fine
6. **Nested arrays:** `[[1, 2], [3, 4]]` - should just work

---

## Example Test Programs

### Basic Arrays
```beeflang
wrangle io

praise ChurchOfBeef():
   prep nums = [1, 2, 3]
   io.preach(nums)        # [1, 2, 3]
   io.preach(nums[0])     # 1
   io.preach(nums[2])     # 3
beef
```

### Array Iteration (Future - needs len())
```beeflang
wrangle io

praise ChurchOfBeef():
   prep nums = [10, 20, 30, 40]
   prep i = 0

   feast while i < len(nums):
      io.preach(nums[i])
      i = i + 1
   beef
beef
```

### Mixed Types
```beeflang
wrangle io

praise ChurchOfBeef():
   prep mixed = [42, "beef", true, [1, 2]]
   io.preach(mixed[0])   # 42
   io.preach(mixed[1])   # beef
   io.preach(mixed[2])   # true
   io.preach(mixed[3])   # [1, 2]
beef
```

---

## Success Criteria

✅ Can parse array literals: `[1, 2, 3]`
✅ Can parse empty arrays: `[]`
✅ Can parse index expressions: `arr[0]`
✅ Arrays can hold mixed types
✅ Can index into arrays to read values
✅ Can index into strings to read characters
✅ Clear error on out-of-bounds access
✅ `len()` builtin works for arrays and strings

**Stretch goals:**
- Index assignment `arr[0] = 99`
- `push()` and `pop()` builtins
- Negative index support
