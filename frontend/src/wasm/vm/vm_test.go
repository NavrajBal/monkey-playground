package vm

import (
	"monkey-wasm/ast"
	"monkey-wasm/compiler"
	"monkey-wasm/lexer"
	"monkey-wasm/object"
	"monkey-wasm/parser"
	"testing"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []struct { input string; expected int64 }{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
	}
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		if err := comp.Compile(program); err != nil { t.Fatalf("compiler error: %s", err) }
		machine := New(comp.Bytecode())
		if err := machine.Run(); err != nil { t.Fatalf("vm error: %s", err) }
		stackElem := machine.LastPoppedStackElem()
		integer, ok := stackElem.(*object.Integer)
		if !ok { t.Fatalf("object is not Integer. got=%T (%+v)", stackElem, stackElem) }
		if integer.Value != tt.expected { t.Fatalf("wrong result. want=%d, got=%d", tt.expected, integer.Value) }
	}
}

func TestClosuresVM(t *testing.T) {
	tests := []struct{ input string; expected int64 }{
		{`let newClosure = fn(a) { fn() { a; }; }; let c = newClosure(99); c();`, 99},
		{`let newAdder = fn(a,b){ fn(c){ a + b + c } }; let adder = newAdder(1,2); adder(8);`, 11},
	}
	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		if err := comp.Compile(program); err != nil { t.Fatalf("compiler error: %s", err) }
		machine := New(comp.Bytecode())
		if err := machine.Run(); err != nil { t.Fatalf("vm error: %s", err) }
		stackElem := machine.LastPoppedStackElem()
		integer, ok := stackElem.(*object.Integer)
		if !ok { t.Fatalf("object is not Integer. got=%T (%+v)", stackElem, stackElem) }
		if integer.Value != tt.expected { t.Fatalf("wrong result. want=%d, got=%d", tt.expected, integer.Value) }
	}
}

func TestRecursiveFunctionsVM(t *testing.T) {
	program := parse(`let countDown = fn(x){ if (x == 0) { return 0; } else { countDown(x - 1); } }; countDown(1);`)
	comp := compiler.New()
	if err := comp.Compile(program); err != nil { t.Fatalf("compiler error: %s", err) }
	machine := New(comp.Bytecode())
	if err := machine.Run(); err != nil { t.Fatalf("vm error: %s", err) }
	stackElem := machine.LastPoppedStackElem()
	integer, ok := stackElem.(*object.Integer)
	if !ok { t.Fatalf("object is not Integer. got=%T (%+v)", stackElem, stackElem) }
	if integer.Value != 0 { t.Fatalf("wrong result. want=%d, got=%d", 0, integer.Value) }
}

func parse(input string) *ast.Program { l := lexer.New(input); p := parser.New(l); return p.ParseProgram() }


