# Beeflang Design Specification

**A meme language honoring the Church of Beef**

## Core Philosophy
Beef-themed keywords with religious overtones. Syntax should be humorous but make intuitive sense where possible.

---

## Agreed Keywords

- `praise` - function declaration
- `beef` - block terminator (end of function/loop/conditional)
- `feast while` - while loop
- `if` / `else` - conditionals
- `cut` - variable declaration
- `serve` - return statement
- `genesis` - entry point (main function)
- `preach` - print/output function

**Current Syntax:**
```
praise genesis():
   cut x = 42
   preach(x)
   
   if x > 0:
      cut y = 5
   beef
   
   feast while x > 0:
      preach(x)
      cut x = x - 1
   beef
beef
```

---

## Operators (Standard)

### Arithmetic
- `+` - addition
- `-` - subtraction
- `*` - multiplication
- `/` - division
- `%` - modulo

### Comparison
- `==` - equal
- `!=` - not equal
- `<` - less than
- `>` - greater than
- `<=` - less than or equal
- `>=` - greater than or equal

### Logical
- `&&` or `and` - logical and
- `||` or `or` - logical or
- `!` or `not` - logical not

### Assignment
- `=` - assignment

---

## Data Types

- **int** - integers
- **bool** - true/false
- **string** - text literals (double-quoted: `"Hello, Beef!"`)

**Type System:** Dynamically typed (inferred from value)

---

## ToDo: Turing Completeness Requirements

### 1. ~~Conditionals (if/else)~~ ✓
- [x] Keywords: `if` / `else`

### 2. ~~Variables~~ ✓
- [x] Variable declaration: `cut`
- [x] Scope rules: Block-level scoping

### 3. ~~Return Statement~~ ✓
- [x] Keyword: `serve`

### 4. ~~Operators~~ ✓
- [x] Standard operators

### 5. ~~Data Types~~ ✓
- [x] int, bool, string (dynamically typed)

### 6. ~~Entry Point~~ ✓
- [x] `genesis` - main function

### 7. ~~Basic I/O~~ ✓
- [x] `preach` - print/output

---

## Language Details

### Scope Rules
- **Block-level scoping**: Variables are scoped to the block they're declared in
- Variables declared in a function/loop/conditional only exist within that block and nested blocks

### Statement Terminators
- **Newline-terminated**: Statements end at newlines
- No semicolons required

### Comments
- **Single-line**: `#` (Python-style)
  ```
  # This is a comment
  cut x = 5  # inline comment
  ```

---

## Complete Example

```beeflang
# Calculate factorial recursively
praise factorial(n):
   if n <= 1:
      serve 1
   beef

   cut result = n * factorial(n - 1)
   serve result
beef

praise genesis():
   preach("Welcome to the Church of Beef!")

   cut num = 5
   cut answer = factorial(num)
   cut message = "The answer is: " + answer
   preach(message)  # outputs: The answer is: 120
beef
```

### 4. Operators
- [ ] Arithmetic: +, -, *, /, %
- [ ] Comparison: ==, !=, <, >, <=, >=
- [ ] Logical: and, or, not
- [ ] Assignment: =

### 5. Data Types
- [ ] What types to support? (int, float, string, bool at minimum?)
- [ ] Type system approach (static/dynamic/inferred?)

### 6. Entry Point
- [ ] Main function keyword/convention
- [ ] Program execution model

### 7. Basic I/O (Optional but useful)
- [ ] Print/output function
- [ ] Input function

---

## Brainstorm 
### Potential Keywords

#### Control Flow
- `sermon` - if statement
- `heresy` - else
- `sacrifice` - break
- `repent` - continue

#### Data/Variables
- `offering` - variable declaration
- `sacred` - constant declaration
- `herd` - array/list
- `congregation` - object/struct
- `scripture` - string type

#### Operators
- `smite` - not/negation
- `unite` - and
- `schism` - or

#### Functions
- `blessing` - return statement

#### Module System
- `summon` - import/include
- `banish` - delete/remove

#### Other
- `genesis` - main/entry point
- `amen` - statement terminator
- `emptiness` or `famine` - null/void value

### Other Ideas

#### Immutable vs Mutable variables
- Mutable: `cut` (v1)
- Immutable: `pack` (v2 - later)