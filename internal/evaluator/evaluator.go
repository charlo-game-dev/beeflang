package evaluator

import (
	"bufio"
	"fmt"
	"os"

	"github.com/elitwilson/beeflang/internal/ast"
	"github.com/elitwilson/beeflang/internal/object"
	"github.com/elitwilson/beeflang/internal/token"
)

// Eval evaluates an AST node and returns the resulting runtime object.
// This is the core of the interpreter - it walks the AST and executes the code.
func Eval(node ast.Node, env *Environment) object.Object {
	switch n := node.(type) {

	// Program: evaluate all statements and return the last result
	case *ast.Program:
		return evalProgram(n, env)

	// Literals: convert AST literals to runtime objects
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	// Identifiers: look up variable in environment
	case *ast.Identifier:
		return evalIdentifier(n, env)

	// Expressions: evaluate recursively
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Token, n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Token, n.Operator, left, right)

	// Statements
	case *ast.VariableDeclaration:
		val := Eval(n.Value, env)
		env.Set(n.Name.Value, val)
		return val

	case *ast.AssignmentStatement:
		return evalAssignmentStatement(n, env)

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.IfStatement:
		return evalIfStatement(n, env)

	case *ast.WhileLoop:
		return evalWhileLoop(n, env)

	case *ast.FunctionDeclaration:
		return evalFunctionDeclaration(n, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(n, env)

	case *ast.FunctionCall:
		return evalFunctionCall(n, env)

	case *ast.WrangleStatement:
		return evalWrangleStatement(n, env)

	case *ast.MemberAccessExpression:
		return evalMemberAccessExpression(n, env)

	// Expression statement: evaluate the expression
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	}

	return nil
}

// evalProgram evaluates all statements in a program and returns the last result
func evalProgram(program *ast.Program, env *Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// Stop evaluation if we hit an error
		if isError(result) {
			return result
		}

		// Stop evaluation if we hit a return statement
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// evalIdentifier looks up a variable in the environment
func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError(node.Token, "identifier not found: %s", node.Value)
	}
	return val
}

// evalPrefixExpression evaluates prefix expressions like -5 or !true
func evalPrefixExpression(tok token.Token, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(tok, right)
	default:
		return newError(tok, "unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperator implements the ! (not) operator
func evalBangOperator(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

// evalMinusPrefixOperator implements the - (negation) operator
func evalMinusPrefixOperator(tok token.Token, right object.Object) object.Object {
	if right.Type() != "INTEGER" {
		return newError(tok, "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalInfixExpression evaluates infix expressions like 5 + 3 or 10 > 5
func evalInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	switch {
	// Integer operations
	case left.Type() == "INTEGER" && right.Type() == "INTEGER":
		return evalIntegerInfixExpression(tok, operator, left, right)

	// String concatenation
	case left.Type() == "STRING" && right.Type() == "STRING":
		return evalStringInfixExpression(tok, operator, left, right)

	// Boolean comparison (using pointer equality optimization)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	// Type mismatch
	case left.Type() != right.Type():
		return newError(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression handles arithmetic and comparison on integers
func evalIntegerInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	// Arithmetic
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}

	// Comparison
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression handles string operations
func evalStringInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// nativeBoolToBooleanObject converts a Go bool to a Boolean object
// Uses singleton TRUE/FALSE for efficiency
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

// evalBlockStatement evaluates a block of statements and returns the last result
// If a return statement is encountered, it stops execution and returns immediately
func evalBlockStatement(block *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// Stop execution if we hit an error
		if isError(result) {
			return result
		}

		// If we hit a return statement, stop executing and bubble it up
		if result != nil && result.Type() == "RETURN_VALUE" {
			return result
		}
	}

	return result
}

// evalIfStatement evaluates an if/else statement
func evalIfStatement(ifStmt *ast.IfStatement, env *Environment) object.Object {
	condition := Eval(ifStmt.Condition, env)

	if isTruthy(condition) {
		return Eval(ifStmt.Consequence, env)
	} else if ifStmt.Alternative != nil {
		return Eval(ifStmt.Alternative, env)
	} else {
		return object.NULL
	}
}

// isTruthy determines if an object is "truthy" for conditionals
// In Beeflang: false and null are falsy, everything else is truthy
func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.FALSE:
		return false
	case object.TRUE:
		return true
	default:
		return true
	}
}

// evalFunctionDeclaration creates a Function object and stores it in the environment
func evalFunctionDeclaration(fn *ast.FunctionDeclaration, env *Environment) object.Object {
	function := &object.Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		Env:        env, // Capture current environment (closure)
	}

	// Store the function in the environment by its name
	env.Set(fn.Name.Value, function)

	return function
}

// evalReturnStatement evaluates a return statement
func evalReturnStatement(stmt *ast.ReturnStatement, env *Environment) object.Object {
	val := Eval(stmt.ReturnValue, env)
	// Wrap in ReturnValue to signal this is an early return
	return &object.ReturnValue{Value: val}
}

// evalFunctionCall evaluates a function call expression
func evalFunctionCall(call *ast.FunctionCall, env *Environment) object.Object {
	// Evaluate the function expression (usually an identifier or member access)
	function := Eval(call.Function, env)
	if isError(function) {
		return function
	}

	// Evaluate all arguments
	args := evalExpressions(call.Arguments, env)
	// Check if any argument evaluation resulted in an error
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	// Check if it's a builtin function
	if builtin, ok := function.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}

	// Check if it's a user-defined function
	fn, ok := function.(*object.Function)
	if !ok {
		// Not a function - error
		return newError(call.Token, "not a function: %s", function.Type())
	}

	// Create new environment for function execution (enclosed by function's closure env)
	fnEnv := object.NewEnclosedEnvironment(fn.Env)

	// Bind parameters to arguments
	for i, param := range fn.Parameters {
		fnEnv.Set(param.Value, args[i])
	}

	// Execute function body
	result := Eval(fn.Body, fnEnv)

	// Propagate errors from function body
	if isError(result) {
		return result
	}

	// Only return a value if there was an explicit "serve" statement
	// Otherwise, functions return NULL (for side-effect-only functions)
	if returnValue, ok := result.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	// No explicit return - function returns NULL
	return object.NULL
}

// evalExpressions evaluates a list of expressions (used for function arguments)
func evalExpressions(exps []ast.Expression, env *Environment) []object.Object {
	result := []object.Object{}

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		result = append(result, evaluated)
	}

	return result
}

// evalAssignmentStatement handles variable reassignment (x = value)
func evalAssignmentStatement(stmt *ast.AssignmentStatement, env *Environment) object.Object {
	val := Eval(stmt.Value, env)
	env.Set(stmt.Name.Value, val)
	return val
}

// evalWhileLoop handles while loops: feast while condition: body beef
func evalWhileLoop(loop *ast.WhileLoop, env *Environment) object.Object {
	var result object.Object = object.NULL

	for {
		condition := Eval(loop.Condition, env)

		if !isTruthy(condition) {
			break
		}

		result = Eval(loop.Body, env)

		// Check for early return from within the loop
		if result != nil && result.Type() == "RETURN_VALUE" {
			return result
		}
	}

	return result
}

func evalWrangleStatement(stmt *ast.WrangleStatement, env *Environment) object.Object {
	// Load module by name
	moduleName := stmt.ModuleName.Value
	mod := loadModule(moduleName)

	// Store module in environment
	env.Set(moduleName, mod)

	return mod
}

func evalMemberAccessExpression(expr *ast.MemberAccessExpression, env *Environment) object.Object {
	// Evaluate the object (left side)
	obj := Eval(expr.Object, env)

	// Check if it's a module
	if mod, ok := obj.(*object.Module); ok {
		member, found := mod.Get(expr.Member.Value)
		if !found {
			return object.NULL
		}
		return member
	}

	return object.NULL
}

// loadModule creates and returns a module by name
// For now, this is hardcoded - later we can make it extensible
func loadModule(name string) *object.Module {
	switch name {
	case "io":
		return createIOModule()
	default:
		// Return empty module for unknown modules
		return &object.Module{
			Name:    name,
			Members: make(map[string]object.Object),
		}
	}
}

func createIOModule() *object.Module {
	mod := &object.Module{
		Name:    "io",
		Members: make(map[string]object.Object),
	}

	// preach - print to stdout with newline
	mod.Set("preach", &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return object.NULL
		},
	})

	// input - read line from stdin
	mod.Set("input", &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// Optional: first argument is prompt
			if len(args) > 0 {
				fmt.Print(args[0].Inspect())
			}

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				return &object.String{Value: scanner.Text()}
			}

			return &object.String{Value: ""}
		},
	})

	return mod
}

// ========================================
// Error Handling Helpers
// ========================================

// newError creates an Error object with a formatted message and location information.
// The token provides line and column numbers for helpful error messages.
//
// Usage: return newError(node.Token, "type mismatch: %s + %s", left.Type(), right.Type())
func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Column:  tok.Column,
		// File is set by main.go when running from a file
	}
}

// isError checks if an object is an Error.
// Used throughout the evaluator to detect and propagate errors up the call stack.
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == "ERROR"
	}
	return false
}
