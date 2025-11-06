package evaluator

import (
	"bytes"

	"github.com/NavrajBal/monkey-lang/object"
)

// Global buffer for WASM output capture
var WasmOutputBuffer *bytes.Buffer

// Initialize the WASM output buffer - called from main
func InitWasmBuffer() {
	if WasmOutputBuffer == nil {
		WasmOutputBuffer = &bytes.Buffer{}
	}
}

var builtins = map[string]*object.Builtin{
	"len":  {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
		switch arg := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	}},
	"puts": {Fn: func(args ...object.Object) object.Object {
		// WASM-compatible puts - write to a global buffer instead of stdout
		// Match VM behavior: print each argument on a new line
		InitWasmBuffer() // Ensure buffer is initialized
		
		for _, arg := range args {
			if arg == nil {
				continue // Skip nil arguments
			}
			WasmOutputBuffer.WriteString(arg.Inspect())
			WasmOutputBuffer.WriteString("\n")
		}
		return NULL
	}},
	"first": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
		if args[0] == nil || args[0].Type() != object.ARRAY_OBJ { return newError("argument to `first` must be ARRAY, got %s", args[0].Type()) }
		arr := args[0].(*object.Array)
		if arr != nil && arr.Elements != nil && len(arr.Elements) > 0 { return arr.Elements[0] }
		return NULL
	}},
	"last": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
		if args[0] == nil || args[0].Type() != object.ARRAY_OBJ { return newError("argument to `last` must be ARRAY, got %s", args[0].Type()) }
		arr := args[0].(*object.Array)
		if arr != nil && arr.Elements != nil {
			length := len(arr.Elements)
			if length > 0 { return arr.Elements[length-1] }
		}
		return NULL
	}},
	"rest": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
		if args[0] == nil || args[0].Type() != object.ARRAY_OBJ { return newError("argument to `rest` must be ARRAY, got %s", args[0].Type()) }
		arr := args[0].(*object.Array)
		if arr != nil && arr.Elements != nil {
			length := len(arr.Elements)
			if length > 1 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:])
				return &object.Array{Elements: newElements}
			}
		}
		// Return empty array if length <= 1
		return &object.Array{Elements: []object.Object{}}
	}},
	"push": {Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 { return newError("wrong number of arguments. got=%d, want=2", len(args)) }
		if args[0] == nil || args[0].Type() != object.ARRAY_OBJ { return newError("argument to `push` must be ARRAY, got %s", args[0].Type()) }
		arr := args[0].(*object.Array)
		if arr != nil && arr.Elements != nil {
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		}
		// Handle nil Elements case
		return &object.Array{Elements: []object.Object{args[1]}}
	}},
}


