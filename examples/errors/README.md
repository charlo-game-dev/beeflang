# Beeflang Error Examples

This directory contains example programs that demonstrate Beeflang's error handling system.

## Running the Examples

Run all error examples:
```bash
./run_error_examples.sh
```

Run a specific example:
```bash
go run ../../main.go type_mismatch.beef
```

## Error Examples

### type_mismatch.beef
Demonstrates type mismatch errors when trying to combine incompatible types.
```
Error at line 10, column 13 - type mismatch: INTEGER + BOOLEAN
```

### undefined_variable.beef
Demonstrates what happens when you reference a variable that doesn't exist.
```
Error at line 10, column 13 - identifier not found: someUndefinedVariable
```

### unknown_operator.beef
Demonstrates invalid operator usage (e.g., adding booleans).
```
Error at line 11, column 17 - unknown operator: BOOLEAN + BOOLEAN
```

### invalid_negation.beef
Demonstrates invalid negation of non-integer types.
```
Error at line 9, column 17 - unknown operator: -BOOLEAN
```

### string_type_mismatch.beef
Demonstrates type mismatch when mixing strings and integers.
```
Error at line 10, column 17 - type mismatch: STRING + INTEGER
```

## Error System Features

All errors include:
- **Line number** - Where the error occurred
- **Column number** - Exact position in the line
- **Clear message** - Explanation of what went wrong
- **Execution stops** - No subsequent code runs after an error

These location details make debugging easy!
