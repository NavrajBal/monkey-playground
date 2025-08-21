package object

import "fmt"

// Builtins lists builtin functions available in the VM runtime
var Builtins = []struct {
    Name    string
    Builtin *Builtin
}{
    {"len", &Builtin{Fn: func(args ...Object) Object {
        if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
        switch arg := args[0].(type) {
        case *Array:
            return &Integer{Value: int64(len(arg.Elements))}
        case *String:
            return &Integer{Value: int64(len(arg.Value))}
        default:
            return newError("argument to `len` not supported, got %s", args[0].Type())
        }
    }}},
    {"puts", &Builtin{Fn: func(args ...Object) Object { for _, arg := range args { fmt.Println(arg.Inspect()) }; return nil }}},
    {"first", &Builtin{Fn: func(args ...Object) Object {
        if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
        if args[0].Type() != ARRAY_OBJ { return newError("argument to `first` must be ARRAY, got %s", args[0].Type()) }
        arr := args[0].(*Array)
        if len(arr.Elements) > 0 { return arr.Elements[0] }
        return nil
    }}},
    {"last", &Builtin{Fn: func(args ...Object) Object {
        if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
        if args[0].Type() != ARRAY_OBJ { return newError("argument to `last` must be ARRAY, got %s", args[0].Type()) }
        arr := args[0].(*Array)
        if l := len(arr.Elements); l > 0 { return arr.Elements[l-1] }
        return nil
    }}},
    {"rest", &Builtin{Fn: func(args ...Object) Object {
        if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=1", len(args)) }
        if args[0].Type() != ARRAY_OBJ { return newError("argument to `rest` must be ARRAY, got %s", args[0].Type()) }
        arr := args[0].(*Array)
        if l := len(arr.Elements); l > 0 { newElements := make([]Object, l-1, l-1); copy(newElements, arr.Elements[1:l]); return &Array{Elements: newElements} }
        return nil
    }}},
    {"push", &Builtin{Fn: func(args ...Object) Object {
        if len(args) != 2 { return newError("wrong number of arguments. got=%d, want=2", len(args)) }
        if args[0].Type() != ARRAY_OBJ { return newError("argument to `push` must be ARRAY, got %s", args[0].Type()) }
        arr := args[0].(*Array)
        l := len(arr.Elements)
        newElements := make([]Object, l+1, l+1); copy(newElements, arr.Elements); newElements[l] = args[1]
        return &Array{Elements: newElements}
    }}},
}

func newError(format string, a ...interface{}) *Error { return &Error{Message: fmt.Sprintf(format, a...)} }

// GetBuiltinByName finds a builtin by name
func GetBuiltinByName(name string) *Builtin {
    for _, def := range Builtins { if def.Name == name { return def.Builtin } }
    return nil
}


