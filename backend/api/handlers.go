package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/NavrajBal/monkey-lang/ast"
	"github.com/NavrajBal/monkey-lang/compiler"
	"github.com/NavrajBal/monkey-lang/evaluator"
	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/object"
	"github.com/NavrajBal/monkey-lang/parser"
	"github.com/NavrajBal/monkey-lang/token"
	"github.com/NavrajBal/monkey-lang/vm"
)

// Request/Response types
type CodeRequest struct {
	Code string `json:"code"`
}

type TokenizeResponse struct {
	Tokens []TokenInfo `json:"tokens"`
	Error  string      `json:"error,omitempty"`
}

type TokenInfo struct {
	Type     string `json:"type"`
	Literal  string `json:"literal"`
	Position int    `json:"position"`
}

type ParseResponse struct {
	AST   interface{} `json:"ast"`
	Error string      `json:"error,omitempty"`
}

type CompileResponse struct {
	Bytecode    []byte        `json:"bytecode"`
	Constants   []interface{} `json:"constants"`
	Instructions string       `json:"instructions"`
	Error       string        `json:"error,omitempty"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ReplResponse struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// TokenizeHandler converts code to tokens
func TokenizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l := lexer.New(req.Code)
	var tokens []TokenInfo
	position := 0

	for {
		tok := l.NextToken()
		if tok.Type == token.EOF {
			break
		}

		tokens = append(tokens, TokenInfo{
			Type:     string(tok.Type),
			Literal:  tok.Literal,
			Position: position,
		})
		position += len(tok.Literal)
	}

	response := TokenizeResponse{Tokens: tokens}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ParseHandler converts code to AST
func ParseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		response := ParseResponse{Error: p.Errors()[0]}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert AST to JSON-serializable format with type information
	astData := ConvertASTToJSON(program)

	response := ParseResponse{AST: astData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompileHandler compiles code to bytecode
func CompileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		response := CompileResponse{Error: p.Errors()[0]}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		response := CompileResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	bytecode := comp.Bytecode()
	
	// Convert constants to JSON-serializable format
	constants := make([]interface{}, len(bytecode.Constants))
	for i, c := range bytecode.Constants {
		constants[i] = c.Inspect()
	}

	response := CompileResponse{
		Bytecode:     bytecode.Instructions,
		Constants:    constants,
		Instructions: bytecode.Instructions.String(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExecuteHandler executes code using the VM
func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		response := ExecuteResponse{Error: p.Errors()[0]}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		response := ExecuteResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Capture stdout during execution
	oldStdout := os.Stdout
	pipeR, pipeW, _ := os.Pipe()
	os.Stdout = pipeW

	// Buffer to capture output
	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, pipeR)
		outputChan <- buf.String()
	}()

	// Run the VM
	machine := vm.New(comp.Bytecode())
	err := machine.Run()

	// Restore stdout
	pipeW.Close()
	os.Stdout = oldStdout
	output := <-outputChan

	if err != nil {
		response := ExecuteResponse{Error: err.Error(), Output: output}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	lastPopped := machine.LastPoppedStackElem()
	result := lastPopped.Inspect()

	response := ExecuteResponse{Result: result, Output: output}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReplHandler provides REPL-like functionality
func ReplHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// For now, just use the evaluator for REPL-like behavior
	l := lexer.New(req.Code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		response := ReplResponse{Error: p.Errors()[0]}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		response := ReplResponse{Result: result.Inspect()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ReplResponse{Result: "null"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConvertASTToJSON converts AST nodes to JSON with type information
func ConvertASTToJSON(node ast.Node) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Get the type name using reflection
	nodeType := reflect.TypeOf(node)
	if nodeType.Kind() == reflect.Ptr {
		nodeType = nodeType.Elem()
	}
	result["type"] = nodeType.Name()
	
	// Add common fields
	result["string"] = node.String()
	
	// Handle specific node types
	switch n := node.(type) {
	case *ast.Program:
		var statements []map[string]interface{}
		for _, stmt := range n.Statements {
			statements = append(statements, ConvertASTToJSON(stmt))
		}
		result["statements"] = statements
		
	case *ast.LetStatement:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		if n.Name != nil {
			result["Name"] = ConvertASTToJSON(n.Name)
		}
		if n.Value != nil {
			result["Value"] = ConvertASTToJSON(n.Value)
		}
		
	case *ast.ReturnStatement:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		if n.ReturnValue != nil {
			result["ReturnValue"] = ConvertASTToJSON(n.ReturnValue)
		}
		
	case *ast.ExpressionStatement:
		if n.Expression != nil {
			result["Expression"] = ConvertASTToJSON(n.Expression)
		}
		
	case *ast.Identifier:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Value"] = n.Value
		
	case *ast.IntegerLiteral:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Value"] = n.Value
		
	case *ast.Boolean:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Value"] = n.Value
		
	case *ast.StringLiteral:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Value"] = n.Value
		
	case *ast.InfixExpression:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Operator"] = n.Operator
		if n.Left != nil {
			result["Left"] = ConvertASTToJSON(n.Left)
		}
		if n.Right != nil {
			result["Right"] = ConvertASTToJSON(n.Right)
		}
		
	case *ast.PrefixExpression:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		result["Operator"] = n.Operator
		if n.Right != nil {
			result["Right"] = ConvertASTToJSON(n.Right)
		}
		
	case *ast.IfExpression:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		if n.Condition != nil {
			result["Condition"] = ConvertASTToJSON(n.Condition)
		}
		if n.Consequence != nil {
			result["Consequence"] = ConvertASTToJSON(n.Consequence)
		}
		if n.Alternative != nil {
			result["Alternative"] = ConvertASTToJSON(n.Alternative)
		}
		
	case *ast.BlockStatement:
		var statements []map[string]interface{}
		for _, stmt := range n.Statements {
			statements = append(statements, ConvertASTToJSON(stmt))
		}
		result["statements"] = statements
		
	case *ast.FunctionLiteral:
		result["Token"] = map[string]interface{}{
			"Type":    string(n.Token.Type),
			"Literal": n.Token.Literal,
		}
		var parameters []map[string]interface{}
		for _, param := range n.Parameters {
			parameters = append(parameters, ConvertASTToJSON(param))
		}
		result["Parameters"] = parameters
		if n.Body != nil {
			result["Body"] = ConvertASTToJSON(n.Body)
		}
		
	case *ast.CallExpression:
		if n.Function != nil {
			result["Function"] = ConvertASTToJSON(n.Function)
		}
		var arguments []map[string]interface{}
		for _, arg := range n.Arguments {
			arguments = append(arguments, ConvertASTToJSON(arg))
		}
		result["Arguments"] = arguments
	}
	
	return result
}
