package object

import (
	"fmt"

	"github.com/elitwilson/beeflang/internal/ast"
)

// Object represents a runtime value in the Beeflang interpreter.
//
// NOTE: "Object" here refers to runtime values (integers, booleans, strings, etc.),
// NOT Object-Oriented Programming. This is standard interpreter terminology - even
// non-OOP languages represent runtime values as "objects."
//
// Every value that exists during program execution implements this interface.
type Object interface {
	Type() string   // Returns the type of the object (e.g., "INTEGER", "BOOLEAN")
	Inspect() string // Returns a string representation for debugging/printing
}

// Integer represents an integer value at runtime.
type Integer struct {
	Value int64
}

func (i *Integer) Type() string {
	return "INTEGER"
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean represents a boolean value at runtime.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() string {
	return "BOOLEAN"
}

func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// String represents a string value at runtime.
type String struct {
	Value string
}

func (s *String) Type() string {
	return "STRING"
}

func (s *String) Inspect() string {
	return s.Value
}

// Null represents the absence of a value.
// Used for functions that don't return anything, uninitialized variables, etc.
type Null struct{}

func (n *Null) Type() string {
	return "NULL"
}

func (n *Null) Inspect() string {
	return "null"
}

// Function represents a function at runtime.
// It stores the function's parameters, body, and the environment where it was defined (closure).
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment // Closure: captures environment where function was defined
}

func (f *Function) Type() string {
	return "FUNCTION"
}

func (f *Function) Inspect() string {
	return "<function>"
}

// ReturnValue wraps a value that's being returned from a function.
// This wrapper allows us to distinguish between a normal evaluation result
// and an early return statement, so we can stop executing and unwind the call stack.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() string {
	return "RETURN_VALUE"
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Environment stores variable bindings (name -> value mappings).
// It supports nested scopes through the `outer` pointer, enabling block-level scoping.
//
// Example:
//   outer := NewEnvironment()
//   outer.Set("x", &Integer{Value: 10})
//
//   inner := NewEnclosedEnvironment(outer)
//   inner.Set("y", &Integer{Value: 20})
//   inner.Get("x")  // finds x in outer scope
//   inner.Get("y")  // finds y in inner scope
type Environment struct {
	store map[string]Object
	outer *Environment // pointer to enclosing (parent) scope
}

// NewEnvironment creates a new environment with no outer scope (global scope).
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new environment enclosed by an outer environment.
// This is used for creating nested scopes (functions, blocks, etc.).
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable from the environment.
// It searches the current scope first, then walks up the outer scopes.
// Returns (value, true) if found, (nil, false) if not found.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		// Not found in current scope, check outer scope
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set stores a variable in the current environment scope.
// This does NOT modify outer scopes - it creates/updates in the current scope only.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Singleton instances used throughout the interpreter for efficiency.
// Instead of creating new objects, we reuse these single instances.
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// Module represents a module/namespace containing functions and values.
// Used for organizing standard library and user modules (e.g., io.preach).
type Module struct {
	Name    string
	Members map[string]Object
}

func (m *Module) Type() string {
	return "MODULE"
}

func (m *Module) Inspect() string {
	return fmt.Sprintf("<module '%s'>", m.Name)
}

// Get retrieves a member from the module by name.
func (m *Module) Get(name string) (Object, bool) {
	obj, ok := m.Members[name]
	return obj, ok
}

// Set stores a member in the module by name.
func (m *Module) Set(name string, val Object) {
	m.Members[name] = val
}

// Builtin represents a built-in function implemented in Go.
// The Fn field is a Go function that takes Object arguments and returns an Object.
type Builtin struct {
	Fn func(args ...Object) Object
}

func (b *Builtin) Type() string {
	return "BUILTIN"
}

func (b *Builtin) Inspect() string {
	return "<builtin>"
}

// Error represents a runtime error in Beeflang.
// When evaluation encounters an error, it returns an Error object
// instead of panicking or returning nil.
//
// Errors are first-class values in Beeflang - they implement the Object
// interface and can be stored, passed around, and inspected. This design
// supports future evolution toward explicit error handling (Result types).
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
