package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Phase 2: Real failing tests - actual assertions

func TestObjectInterface(t *testing.T) {
	// Verify all types implement the Object interface
	var _ Object = &Integer{}
	var _ Object = &Boolean{}
	var _ Object = &String{}
	var _ Object = &Null{}
	var _ Object = &Module{}
	var _ Object = &Builtin{}
}

func TestIntegerTypeAndInspect(t *testing.T) {
	integer := &Integer{Value: 42}

	assert.Equal(t, "INTEGER", integer.Type())
	assert.Equal(t, "42", integer.Inspect())
}

func TestBooleanTypeAndInspect(t *testing.T) {
	trueVal := &Boolean{Value: true}
	falseVal := &Boolean{Value: false}

	assert.Equal(t, "BOOLEAN", trueVal.Type())
	assert.Equal(t, "true", trueVal.Inspect())

	assert.Equal(t, "BOOLEAN", falseVal.Type())
	assert.Equal(t, "false", falseVal.Inspect())
}

func TestStringTypeAndInspect(t *testing.T) {
	str := &String{Value: "Hello, Beef!"}

	assert.Equal(t, "STRING", str.Type())
	assert.Equal(t, "Hello, Beef!", str.Inspect())
}

func TestNullTypeAndInspect(t *testing.T) {
	null := &Null{}

	assert.Equal(t, "NULL", null.Type())
	assert.Equal(t, "null", null.Inspect())
}

func TestIntegerValue(t *testing.T) {
	tests := []struct {
		value    int64
		expected int64
	}{
		{42, 42},
		{0, 0},
		{-100, -100},
		{9999, 9999},
	}

	for _, tt := range tests {
		integer := &Integer{Value: tt.value}
		assert.Equal(t, tt.expected, integer.Value)
	}
}

func TestBooleanValue(t *testing.T) {
	trueVal := &Boolean{Value: true}
	falseVal := &Boolean{Value: false}

	assert.True(t, trueVal.Value)
	assert.False(t, falseVal.Value)
}

func TestStringValue(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"Hello, Beef!", "Hello, Beef!"},
		{"", ""},
		{"Praise the Beef!", "Praise the Beef!"},
	}

	for _, tt := range tests {
		str := &String{Value: tt.value}
		assert.Equal(t, tt.expected, str.Value)
	}
}

func TestNullIsUnique(t *testing.T) {
	// Create a global NULL instance for efficiency
	// All null values should reference the same instance
	null1 := NULL
	null2 := NULL

	assert.Equal(t, null1, null2, "NULL should be a singleton")
}

// ========================================
// Module System Object Tests
// ========================================

func TestModuleTypeAndInspect(t *testing.T) {
	mod := &Module{
		Name:    "io",
		Members: make(map[string]Object),
	}

	assert.Equal(t, "MODULE", mod.Type())
	assert.Equal(t, "<module 'io'>", mod.Inspect())
}

func TestModuleGetMember(t *testing.T) {
	mod := &Module{
		Name:    "io",
		Members: make(map[string]Object),
	}

	// Add a member
	testVal := &Integer{Value: 42}
	mod.Members["test"] = testVal

	// Get should retrieve it
	result, ok := mod.Get("test")
	assert.True(t, ok, "Get should return true when member exists")
	assert.Equal(t, testVal, result)

	// Get non-existent member
	_, ok = mod.Get("nonexistent")
	assert.False(t, ok, "Get should return false when member doesn't exist")
}

func TestModuleSetMember(t *testing.T) {
	mod := &Module{
		Name:    "io",
		Members: make(map[string]Object),
	}

	// Set a member
	testVal := &Integer{Value: 42}
	mod.Set("test", testVal)

	// Verify it was stored
	result, ok := mod.Members["test"]
	assert.True(t, ok)
	assert.Equal(t, testVal, result)
}

func TestBuiltinTypeAndInspect(t *testing.T) {
	builtin := &Builtin{
		Fn: func(args ...Object) Object {
			return NULL
		},
	}

	assert.Equal(t, "BUILTIN", builtin.Type())
	assert.Equal(t, "<builtin>", builtin.Inspect())
}

func TestBuiltinFunction(t *testing.T) {
	// Create a builtin that returns the first argument
	builtin := &Builtin{
		Fn: func(args ...Object) Object {
			if len(args) > 0 {
				return args[0]
			}
			return NULL
		},
	}

	// Call it
	arg := &Integer{Value: 42}
	result := builtin.Fn(arg)

	assert.Equal(t, arg, result)
}

// ========================================
// Error Object Tests
// ========================================

// Phase 2: Real failing tests
func TestErrorObjectType(t *testing.T) {
	err := &Error{Message: "something went wrong"}
	assert.Equal(t, "ERROR", err.Type())
}

func TestErrorObjectInspectWithoutLocation(t *testing.T) {
	err := &Error{
		Message: "something went wrong",
		Line:    0,
		Column:  0,
		File:    "",
	}
	assert.Equal(t, "Error: something went wrong", err.Inspect())
}

func TestErrorObjectInspectWithLineColumn(t *testing.T) {
	err := &Error{
		Message: "index out of bounds",
		Line:    4,
		Column:  10,
		File:    "",
	}
	assert.Equal(t, "Error at line 4, column 10 - index out of bounds", err.Inspect())
}

func TestErrorObjectInspectWithFile(t *testing.T) {
	err := &Error{
		Message: "type mismatch",
		Line:    12,
		Column:  5,
		File:    "examples/test.beef",
	}
	assert.Equal(t, "Error at examples/test.beef:12:5 - type mismatch", err.Inspect())
}

func TestErrorImplementsObjectInterface(t *testing.T) {
	var _ Object = &Error{}
}
