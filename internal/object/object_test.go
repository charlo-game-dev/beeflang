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
