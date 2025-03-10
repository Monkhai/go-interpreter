package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: lenFn},
	"print": {Fn: printFn},
	"first": {Fn: firstFn},
	"last":  {Fn: lastFn},
	"rest":  {Fn: restFn},
	"push":  {Fn: push},
}

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	// --------------------------------
	// --------------------------------
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	// --------------------------------
	// --------------------------------
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	// --------------------------------
	// --------------------------------
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}
	// --------------------------------
	// --------------------------------
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	// --------------------------------
	// --------------------------------
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	// --------------------------------
	// --------------------------------
	// Expressions
	case *ast.NullExpression:
		return &object.Null{}
	// --------------------------------
	// --------------------------------
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	// --------------------------------
	// --------------------------------
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	// --------------------------------
	// --------------------------------
	case *ast.Boolean:
		return nativeBoolToObj(node.Value)
	// --------------------------------
	// --------------------------------
	case *ast.Identifier:
		return evalIdentifier(node, env)
	// --------------------------------
	// --------------------------------
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)
	// --------------------------------
	// --------------------------------
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)
	// --------------------------------
	// --------------------------------
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	// --------------------------------
	// --------------------------------
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	// --------------------------------
	// --------------------------------
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	// --------------------------------
	// --------------------------------
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexEpxression(left, index)
	// --------------------------------
	// --------------------------------
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object
	for _, statement := range statements {
		obj = Eval(statement, env)

		switch obj := obj.(type) {
		case *object.ReturnValue:
			return obj.Value
		case *object.Error:
			return obj
		}
	}
	return obj
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, e := range expressions {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}

	return results
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object

	for _, s := range block.Statements {
		obj = Eval(s, env)

		if obj != nil {
			objType := obj.Type()
			if objType == object.RETURN_VALUE_OBJ || objType == object.ERROR_OBJ {
				return obj
			}
		}
	}

	return obj
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}

}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToObj(left == right)
	case operator == "!=":
		return nativeBoolToObj(left != right)
	case right.Type() != left.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToObj(leftVal > rightVal)
	case "<":
		return nativeBoolToObj(leftVal < rightVal)
	case "==":
		return nativeBoolToObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToObj(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if val, ok := builtins[node.Value]; ok {
		return val
	}
	return newError("identifier not found: %s", node.Value)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		fnEnv := extendFuncEnv(function, args)
		evaluated := Eval(function.Body, fnEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError("not a function: %s", function.Type())

	}
}

func extendFuncEnv(fn *object.Function, args []object.Object) *object.Environment {
	enclosedEnv := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		arg := args[i]
		enclosedEnv.Set(param.Value, arg)
	}
	return enclosedEnv
}

func evalIndexEpxression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpresion(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexEpxpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpresion(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexEpxpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	hashableKey, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	hashPair, ok := hashObject.Pairs[hashableKey.HashKey()]
	if !ok {
		return NULL
	}
	return hashPair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

// ================================================
// ================HELPER FUNCTIONS================
// ================================================
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func nativeBoolToObj(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Messgae: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}
