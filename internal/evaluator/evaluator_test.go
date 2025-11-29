package evaluator

import (
	"testing"

	"github.com/elitwilson/beeflang/internal/lexer"
	"github.com/elitwilson/beeflang/internal/object"
	"github.com/elitwilson/beeflang/internal/parser"
	"github.com/stretchr/testify/assert"
)

// Phase 2: Real failing tests - basic expressions

// Helper function to parse and evaluate Beeflang source code
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	return Eval(program, env)
}

func TestEvalIntegerLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"42", 42},
		{"0", 0},
		{"999", 999},
		{"-5", -5},
		{"-100", -100},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value")

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer object")
		assert.Equal(t, tt.expected, integer.Value)
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value")

		boolean, ok := result.(*object.Boolean)
		assert.True(t, ok, "Result should be a Boolean object")
		assert.Equal(t, tt.expected, boolean.Value)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello, Beef!"`, "Hello, Beef!"},
		{`"Praise the Beef!"`, "Praise the Beef!"},
		{`""`, ""},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value")

		str, ok := result.(*object.String)
		assert.True(t, ok, "Result should be a String object")
		assert.Equal(t, tt.expected, str.Value)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Integer negation
		{"-5", int64(-5)},
		{"-42", int64(-42)},
		{"--10", int64(10)},

		// Boolean negation
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value for input: %s", tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			integer, ok := result.(*object.Integer)
			assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
			assert.Equal(t, expected, integer.Value, "Input: %s", tt.input)
		case bool:
			boolean, ok := result.(*object.Boolean)
			assert.True(t, ok, "Result should be a Boolean for input: %s", tt.input)
			assert.Equal(t, expected, boolean.Value, "Input: %s", tt.input)
		}
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Integer arithmetic
		{"5 + 5", int64(10)},
		{"10 - 5", int64(5)},
		{"2 * 3", int64(6)},
		{"10 / 2", int64(5)},
		{"10 % 3", int64(1)},

		// Integer comparisons
		{"5 < 10", true},
		{"10 > 5", true},
		{"5 == 5", true},
		{"5 != 10", true},
		{"5 <= 5", true},
		{"10 >= 5", true},

		// Boolean operations
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},

		// String concatenation
		{`"Hello" + " " + "Beef"`, "Hello Beef"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value for input: %s", tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			integer, ok := result.(*object.Integer)
			assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
			assert.Equal(t, expected, integer.Value, "Input: %s", tt.input)
		case bool:
			boolean, ok := result.(*object.Boolean)
			assert.True(t, ok, "Result should be a Boolean for input: %s", tt.input)
			assert.Equal(t, expected, boolean.Value, "Input: %s", tt.input)
		case string:
			str, ok := result.(*object.String)
			assert.True(t, ok, "Result should be a String for input: %s", tt.input)
			assert.Equal(t, expected, str.Value, "Input: %s", tt.input)
		}
	}
}

// Phase 1.5: Real failing tests for variables

func TestEvalVariableDeclaration(t *testing.T) {
	input := `
prep x = 42
prep y = x + 8
y
`
	result := testEval(input)
	assert.NotNil(t, result)

	integer, ok := result.(*object.Integer)
	assert.True(t, ok, "Result should be an Integer")
	assert.Equal(t, int64(50), integer.Value)
}

func TestEvalIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"prep a = 5; a", 5},
		{"prep a = 5 * 5; a", 25},
		{"prep a = 5; prep b = a; b", 5},
		{"prep a = 5; prep b = a; prep c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Input: %s", tt.input)

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
		assert.Equal(t, tt.expected, integer.Value, "Input: %s", tt.input)
	}
}

// Phase 2: Control Flow - Real failing tests

func TestEvalBlockStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Block should evaluate all statements and return the last value
		{`
prep x = 5
prep y = 10
x + y
`, 15},
		// Single statement in a block
		{`
42
`, 42},
		// Multiple statements, last one is the result
		{`
prep a = 1
prep b = 2
prep c = 3
c
`, 3},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Input: %s", tt.input)

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
		assert.Equal(t, tt.expected, integer.Value, "Input: %s", tt.input)
	}
}

func TestEvalIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{} // can be int64, bool, or nil (for NULL)
	}{
		// Basic if with true condition
		{"if true: 10 beef", int64(10)},
		// Basic if with false condition (should return NULL)
		{"if false: 10 beef", nil},
		// If with expression condition
		{"if 1 < 2: 10 beef", int64(10)},
		{"if 1 > 2: 10 beef", nil},
		// If-else with true condition
		{"if 1 < 2: 10 else: 20 beef", int64(10)},
		// If-else with false condition
		{"if 1 > 2: 10 else: 20 beef", int64(20)},
		// Boolean result
		{"if true: true else: false beef", true},
		// Nested expressions in consequence
		{"if true: 5 + 5 beef", int64(10)},
		// Nested expressions in alternative
		{"if false: 5 else: 10 + 10 beef", int64(20)},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Eval should return a value for input: %s", tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			integer, ok := result.(*object.Integer)
			assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
			assert.Equal(t, expected, integer.Value, "Input: %s", tt.input)
		case bool:
			boolean, ok := result.(*object.Boolean)
			assert.True(t, ok, "Result should be a Boolean for input: %s", tt.input)
			assert.Equal(t, expected, boolean.Value, "Input: %s", tt.input)
		case nil:
			null, ok := result.(*object.Null)
			assert.True(t, ok, "Result should be NULL for input: %s", tt.input)
			assert.NotNil(t, null, "NULL object should not be nil pointer")
		}
	}
}

// Phase 3: Functions - Real failing tests

func TestEvalFunctionDeclaration(t *testing.T) {
	input := `
praise add(x, y):
   serve x + y
beef
`
	result := testEval(input)
	assert.NotNil(t, result)

	// Function declaration should return a Function object
	fn, ok := result.(*object.Function)
	assert.True(t, ok, "Result should be a Function object")
	assert.Len(t, fn.Parameters, 2, "Function should have 2 parameters")
	assert.Equal(t, "x", fn.Parameters[0].Value)
	assert.Equal(t, "y", fn.Parameters[1].Value)
	assert.NotNil(t, fn.Body, "Function should have a body")
}

func TestEvalFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Simple function call with explicit return
		{`
praise add(x, y):
   serve x + y
beef
add(5, 3)
`, 8},
		// Function with single parameter
		{`
praise double(x):
   serve x * 2
beef
double(4)
`, 8},
		// Function with no parameters
		{`
praise fortytwo():
   serve 42
beef
fortytwo()
`, 42},
		// Multiple calls to same function
		{`
praise add(x, y):
   serve x + y
beef
prep a = add(1, 2)
prep b = add(3, 4)
a + b
`, 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Input: %s", tt.input)

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
		assert.Equal(t, tt.expected, integer.Value, "Input: %s", tt.input)
	}
}

func TestEvalFunctionWithoutReturn(t *testing.T) {
	// Function without explicit serve should return NULL
	input := `
praise noReturn(x, y):
   x + y
beef
noReturn(5, 3)
`
	result := testEval(input)
	assert.NotNil(t, result)

	null, ok := result.(*object.Null)
	assert.True(t, ok, "Function without serve should return NULL")
	assert.NotNil(t, null)
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Simple return
		{`
praise getValue():
   serve 42
beef
getValue()
`, 42},
		// Early return
		{`
praise earlyReturn():
   serve 10
   prep x = 99
   x
beef
earlyReturn()
`, 10},
		// Return with expression
		{`
praise calculate():
   serve 5 + 5
beef
calculate()
`, 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Input: %s", tt.input)

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
		assert.Equal(t, tt.expected, integer.Value, "Input: %s", tt.input)
	}
}

// Phase 4: Loops - Real failing tests

func TestEvalWhileLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// Basic countdown loop
		{`
prep counter = 5
feast while counter > 0:
   counter = counter - 1
beef
counter
`, 0},
		// Loop with accumulator
		{`
prep sum = 0
prep i = 1
feast while i <= 5:
   sum = sum + i
   i = i + 1
beef
sum
`, 15}, // 1+2+3+4+5 = 15
		// Loop that doesn't execute
		{`
prep x = 0
feast while x > 10:
   x = x + 1
beef
x
`, 0},
		// Nested variable mutation
		{`
prep result = 1
prep count = 5
feast while count > 0:
   result = result * 2
   count = count - 1
beef
result
`, 32}, // 1 * 2^5 = 32
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		assert.NotNil(t, result, "Input: %s", tt.input)

		integer, ok := result.(*object.Integer)
		assert.True(t, ok, "Result should be an Integer for input: %s", tt.input)
		assert.Equal(t, tt.expected, integer.Value, "Input: %s", tt.input)
	}
}

// ========================================
// Error Infrastructure Tests
// ========================================

// Phase 2: Real failing tests
func TestIsErrorHelper(t *testing.T) {
	// isError should return true for Error objects
	err := &object.Error{Message: "test error"}
	assert.True(t, isError(err))

	// isError should return false for other objects
	integer := &object.Integer{Value: 42}
	assert.False(t, isError(integer))

	boolean := &object.Boolean{Value: true}
	assert.False(t, isError(boolean))

	// isError should return false for nil
	assert.False(t, isError(nil))
}

func TestErrorPropagation(t *testing.T) {
	// Errors should stop evaluation and propagate up
	input := `
5 + true
10
`
	result := testEval(input)

	errObj, ok := result.(*object.Error)
	assert.True(t, ok, "Expected error object")
	assert.Contains(t, errObj.Message, "type mismatch")
	// The second statement (10) should NOT be evaluated
}

func TestTypeErrorMessages(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch"},
		{"5 + \"hello\"", "type mismatch"},
		{"true + false", "unknown operator"},
		{"-true", "unknown operator"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		errObj, ok := result.(*object.Error)
		assert.True(t, ok, "Expected error for input: %s", tt.input)
		assert.Contains(t, errObj.Message, tt.expectedMessage, "Input: %s", tt.input)
		assert.Greater(t, errObj.Line, 0, "Error should have line number")
	}
}

func TestUndefinedVariableError(t *testing.T) {
	input := "foobar"
	result := testEval(input)

	errObj, ok := result.(*object.Error)
	assert.True(t, ok, "Expected error object")
	assert.Contains(t, errObj.Message, "identifier not found")
	assert.Contains(t, errObj.Message, "foobar")
}

func TestUnknownOperatorError(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"-true"},
		{"true + false"},
		{"5; true + false; 5"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		errObj, ok := result.(*object.Error)
		assert.True(t, ok, "Expected error for input: %s", tt.input)
		assert.Contains(t, errObj.Message, "unknown operator", "Input: %s", tt.input)
	}
}

func TestErrorStopsEvaluation(t *testing.T) {
	// When an error occurs, subsequent statements should not be evaluated
	input := `
prep x = 5
prep y = x + true
prep z = 10
z
`
	result := testEval(input)

	errObj, ok := result.(*object.Error)
	assert.True(t, ok, "Expected error object")
	// The error should be about the type mismatch, not about z being undefined
	assert.Contains(t, errObj.Message, "type mismatch")
}

func TestErrorStopsBlockEvaluation(t *testing.T) {
	// Errors should stop evaluation within a block (like function bodies)
	input := `
praise testFunc():
   prep x = 5
   prep y = x + true
   prep z = 10
   z
beef

testFunc()
`
	result := testEval(input)

	errObj, ok := result.(*object.Error)
	assert.True(t, ok, "Expected error object")
	assert.Contains(t, errObj.Message, "type mismatch")
}

func TestErrorPropagatesFromFunctionCall(t *testing.T) {
	// Errors in function bodies should propagate to the caller
	input := `
praise badFunc():
   serve 5 + true
beef

prep result = badFunc()
`
	result := testEval(input)

	errObj, ok := result.(*object.Error)
	assert.True(t, ok, "Expected error object")
	assert.Contains(t, errObj.Message, "type mismatch")
}
