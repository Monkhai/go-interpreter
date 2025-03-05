package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"-50 + 100 + -50", 0},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestStringObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar"`, "foobar"},
		{`"hello world"`, "hello world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 > 2", false},
		{"2 > 1", true},
		{"2 < 1", false},
		{"1 < 2", true},
		{"1 == 2", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"true == false", false},
		{"true == true", true},
		{"true != true", false},
		{"true != false", true},
		{"null != null", false},
		{"null == null", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		expected, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(expected))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;	
			}
				return 1;
		}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true;", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; };", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if (10 > 1) { 
			if (10 > 1) {
				return true + false;
				}
				return 1;
				};`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar;", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Messgae != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Messgae)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a= 5; a;", 5},
		{"let a= 5 * 5; a;", 25},
		{"let a= 5; let b = a; b;", 5},
		{"let a= 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not object.Function, got=%T", evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("fn.Parameters has unexpected length of params. got=%d. wanted 1", len(fn.Parameters))
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("fn.Parameters[0].String() is not x. got=%s", fn.Parameters[0].String())
	}
	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("fn.Body.String() is not '(x + 2). got %q", fn.Body.String())
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5) );", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not object.Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object.Value has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	if result.Type() != object.INTEGER_OBJ {
		t.Errorf("object.Type has wrong type. got=%s, want=%s", result.Type(), object.INTEGER_OBJ)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not object.String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object.Value has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	if result.Type() != object.STRING_OBJ {
		t.Errorf("object.Type has wrong type. got=%s, want=%s", result.Type(), object.INTEGER_OBJ)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not object.Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object.Value has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	if result.Type() != object.BOOLEAN_OBJ {
		t.Errorf("object.Type has wrong type. got=%s, want=%s", result.Type(), object.BOOLEAN_OBJ)
		return false
	}
	return true
}
