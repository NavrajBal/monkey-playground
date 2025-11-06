//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"monkey-wasm/api"
	"monkey-wasm/evaluator"

	"github.com/NavrajBal/monkey-lang/compiler"
	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/object"
	"github.com/NavrajBal/monkey-lang/parser"
)

// TokenInfo represents a token for WASM
type TokenInfo struct {
	Type     string `json:"type"`
	Literal  string `json:"literal"`
	Position int    `json:"position"`
}

// Note: Output buffer is now in evaluator.WasmOutputBuffer


// WASM function to tokenize Monkey code
func tokenize(this js.Value, args []js.Value) (result any) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panics gracefully and return error response
			fmt.Printf("WASM tokenize panic: %v\n", r)
			result = js.ValueOf(map[string]any{
				"error": fmt.Sprintf("WASM panic: %v", r),
			})
		}
	}()
	
	if len(args) != 1 {
		return js.ValueOf(map[string]any{
			"error": "tokenize requires exactly 1 argument (code string)",
		})
	}

	code := args[0].String()
	if code == "" {
		jsonBytes, _ := json.Marshal(map[string]any{
			"tokens": []TokenInfo{},
		})
		jsonStr := string(jsonBytes)
		jsonParser := js.Global().Get("JSON")
		return jsonParser.Call("parse", jsonStr)
	}
	
	l := lexer.New(code)
	
	var tokens []TokenInfo
	position := 0
	
	for {
		tok := l.NextToken()
		if tok.Type == "EOF" {
			break
		}
		
		tokens = append(tokens, TokenInfo{
			Type:     string(tok.Type),
			Literal:  tok.Literal,
			Position: position,
		})
		position += len(tok.Literal)
	}

	// Use JSON encoding for proper serialization
	jsonBytes, err := json.Marshal(map[string]any{
		"tokens": tokens,
	})
	if err != nil {
		return js.ValueOf(map[string]any{
			"error": fmt.Sprintf("Failed to marshal tokens: %v", err),
		})
	}

	// Parse JSON string to JavaScript object
	jsonStr := string(jsonBytes)
	jsonParser := js.Global().Get("JSON")
	parsed := jsonParser.Call("parse", jsonStr)

	return parsed
}

// WASM function to parse Monkey code to AST
func parseAST(this js.Value, args []js.Value) (result any) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panics gracefully and return error response
			fmt.Printf("WASM parseAST panic: %v\n", r)
			result = js.ValueOf(map[string]any{
				"error": fmt.Sprintf("WASM panic: %v", r),
			})
		}
	}()
	
	if len(args) != 1 {
		return js.ValueOf(map[string]any{
			"error": "parseAST requires exactly 1 argument (code string)",
		})
	}

	code := args[0].String()
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return js.ValueOf(map[string]any{
			"error": p.Errors()[0],
		})
	}

	// Convert AST to JSON-serializable format
	astData := api.ConvertASTToJSON(program)

	// Use JSON encoding to properly serialize nested structures
	// js.ValueOf doesn't handle deeply nested maps correctly, so we use JSON
	jsonBytes, err := json.Marshal(map[string]any{
		"ast": astData,
	})
	if err != nil {
		return js.ValueOf(map[string]any{
			"error": fmt.Sprintf("Failed to marshal AST: %v", err),
		})
	}

	// Parse JSON string to JavaScript object using JSON.parse
	jsonStr := string(jsonBytes)
	jsonParser := js.Global().Get("JSON")
	parsed := jsonParser.Call("parse", jsonStr)

	return parsed
}

// WASM function to compile Monkey code
func compile(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return js.ValueOf(map[string]any{
			"error": "compile requires exactly 1 argument (code string)",
		})
	}

	code := args[0].String()
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return js.ValueOf(map[string]any{
			"error": p.Errors()[0],
		})
	}

	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		return js.ValueOf(map[string]any{
			"error": err.Error(),
		})
	}

	// Get bytecode instructions
	bytecode := comp.Bytecode()
	instructions := bytecode.Instructions.String()

	result := map[string]any{
		"instructions": instructions,
		"constants":    len(bytecode.Constants),
	}

	return js.ValueOf(result)
}

// WASM function to execute Monkey code
func execute(this js.Value, args []js.Value) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panics gracefully and return error response
			fmt.Printf("WASM execute panic: %v\n", r)
			result = js.ValueOf(map[string]any{
				"error": fmt.Sprintf("WASM panic: %v", r),
			})
		}
	}()
	
	if len(args) != 1 {
		return js.ValueOf(map[string]any{
			"error": "execute requires exactly 1 argument (code string)",
		})
	}

	code := args[0].String()
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return js.ValueOf(map[string]any{
			"error": p.Errors()[0],
		})
	}

	// Clear the global output buffer before execution
	evaluator.InitWasmBuffer()
	evaluator.WasmOutputBuffer.Reset()

	// Use the evaluator with default environment (puts builtin is now WASM-compatible)
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	
	if evaluated != nil {
		if errorObj, ok := evaluated.(*object.Error); ok {
			return js.ValueOf(map[string]any{
				"error": errorObj.Message,
			})
		}
	}

	// Get the captured output
	var capturedOutput string
	if evaluator.WasmOutputBuffer != nil {
		capturedOutput = evaluator.WasmOutputBuffer.String()
	}
	
	var resultStr string
	if evaluated != nil {
		resultStr = evaluated.Inspect()
	} else {
		resultStr = "null"
	}
	
	responseData := map[string]any{
		"result": resultStr,
		"output": capturedOutput,
	}

	fmt.Printf("WASM execute returning: %+v\n", responseData)
	return js.ValueOf(responseData)
}

// WASM function for REPL-style evaluation
func repl(this js.Value, args []js.Value) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panics gracefully and return error response
			fmt.Printf("WASM repl panic: %v\n", r)
			result = js.ValueOf(map[string]any{
				"error": fmt.Sprintf("WASM panic: %v", r),
			})
		}
	}()
	
	if len(args) != 1 {
		return js.ValueOf(map[string]any{
			"error": "repl requires exactly 1 argument (code string)",
		})
	}

	code := args[0].String()
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return js.ValueOf(map[string]any{
			"error": p.Errors()[0],
		})
	}

	// Use evaluator for REPL (more interactive)
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		result := map[string]any{
			"result": evaluated.Inspect(),
		}
		return js.ValueOf(result)
	}

	return js.ValueOf(map[string]any{
		"result": "null",
	})
}

func main() {
	// Initialize WASM output buffer
	evaluator.InitWasmBuffer()
	
	// Create a channel to keep the program running
	done := make(chan struct{})

	// Register WASM functions that won't be garbage collected
	tokenizeFunc := js.FuncOf(tokenize)
	parseFunc := js.FuncOf(parseAST)
	compileFunc := js.FuncOf(compile)
	executeFunc := js.FuncOf(execute)
	replFunc := js.FuncOf(repl)

	js.Global().Set("monkeyTokenize", tokenizeFunc)
	js.Global().Set("monkeyParseAST", parseFunc)
	js.Global().Set("monkeyCompile", compileFunc)
	js.Global().Set("monkeyExecute", executeFunc)
	js.Global().Set("monkeyRepl", replFunc)

	// Signal that WASM is ready
	js.Global().Set("monkeyWasmReady", js.ValueOf(true))

	// Add a cleanup function that can be called from JS
	cleanupFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		tokenizeFunc.Release()
		parseFunc.Release()
		compileFunc.Release()
		executeFunc.Release()
		replFunc.Release()
		close(done)
		return nil
	})
	js.Global().Set("monkeyCleanup", cleanupFunc)

	// Keep the program running until cleanup is called
	<-done
}