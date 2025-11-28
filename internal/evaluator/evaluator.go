package evaluator

import (
	"github.com/elitwilson/beeflang/internal/ast"
	"github.com/elitwilson/beeflang/internal/object"
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
		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		right := Eval(n.Right, env)
		return evalInfixExpression(n.Operator, left, right)

	// Statements
	case *ast.VariableDeclaration:
		val := Eval(n.Value, env)
		env.Set(n.Name.Value, val)
		return val

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.IfStatement:
		return evalIfStatement(n, env)

	case *ast.FunctionDeclaration:
		return evalFunctionDeclaration(n, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(n, env)

	case *ast.FunctionCall:
		return evalFunctionCall(n, env)

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
	}

	return result
}

// evalIdentifier looks up a variable in the environment
func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		// Variable not found - return null for now (later we'll add error handling)
		return object.NULL
	}
	return val
}

// evalPrefixExpression evaluates prefix expressions like -5 or !true
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return object.NULL
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
func evalMinusPrefixOperator(right object.Object) object.Object {
	if right.Type() != "INTEGER" {
		return object.NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalInfixExpression evaluates infix expressions like 5 + 3 or 10 > 5
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// Integer operations
	case left.Type() == "INTEGER" && right.Type() == "INTEGER":
		return evalIntegerInfixExpression(operator, left, right)

	// String concatenation
	case left.Type() == "STRING" && right.Type() == "STRING":
		return evalStringInfixExpression(operator, left, right)

	// Boolean comparison (using pointer equality optimization)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return object.NULL
	}
}

// evalIntegerInfixExpression handles arithmetic and comparison on integers
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
		return object.NULL
	}
}

// evalStringInfixExpression handles string operations
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return object.NULL
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
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
	// Evaluate the function expression (usually an identifier)
	function := Eval(call.Function, env)

	fn, ok := function.(*object.Function)
	if !ok {
		// Not a function - error
		return object.NULL
	}

	// Evaluate all arguments
	args := evalExpressions(call.Arguments, env)

	// Create new environment for function execution (enclosed by function's closure env)
	fnEnv := object.NewEnclosedEnvironment(fn.Env)

	// Bind parameters to arguments
	for i, param := range fn.Parameters {
		fnEnv.Set(param.Value, args[i])
	}

	// Execute function body
	result := Eval(fn.Body, fnEnv)

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
