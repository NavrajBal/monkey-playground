package handler

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/NavrajBal/monkey-lang/ast"
	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/parser"
)

type CodeRequest struct {
	Code string `json:"code"`
}

type ParseResponse struct {
	AST   interface{} `json:"ast"`
	Error string      `json:"error,omitempty"`
}

// Handler is the main Vercel function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

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

	// Convert AST to JSON-serializable format
	astData := ConvertASTToJSON(program)

	response := ParseResponse{AST: astData}
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
