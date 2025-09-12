package evaluator

import (
	"testing"

	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/object"
	"github.com/NavrajBal/monkey-lang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct { input string; expected int64 }{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
	}
	for _, tt := range tests { evaluated := testEval(tt.input); testIntegerObject(t, evaluated, tt.expected) }
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct { input string; expected bool }{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 != 2", true},
	}
	for _, tt := range tests { evaluated := testEval(tt.input); testBooleanObject(t, evaluated, tt.expected) }
}

func TestBangOperator(t *testing.T) {
	tests := []struct { input string; expected bool }{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
	}
	for _, tt := range tests { evaluated := testEval(tt.input); testBooleanObject(t, evaluated, tt.expected) }
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct { input string; expected interface{} }{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if v, ok := tt.expected.(int); ok { testIntegerObject(t, evaluated, int64(v)) } else { testNullObject(t, evaluated) }
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct { input string; expected int64 }{
		{"return 10;", 10},
		{"return 2 * 5;", 10},
	}
	for _, tt := range tests { evaluated := testEval(tt.input); testIntegerObject(t, evaluated, tt.expected) }
}

func TestLetStatements(t *testing.T) {
	tests := []struct { input string; expected int64 }{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
	}
	for _, tt := range tests { testIntegerObject(t, testEval(tt.input), tt.expected) }
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	if _, ok := evaluated.(*object.Function); !ok { t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated) }
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct { input string; expected int64 }{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
	}
	for _, tt := range tests { testIntegerObject(t, testEval(tt.input), tt.expected) }
}

// helpers
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok { t.Errorf("object is not Integer. got=%T (%+v)", obj, obj); return false }
	if result.Value != expected { t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected); return false }
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok { t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj); return false }
	if result.Value != expected { t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected); return false }
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL { t.Errorf("object is not NULL. got=%T (%+v)", obj, obj); return false }
	return true
}


