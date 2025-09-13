package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NavrajBal/monkey-lang/compiler"
	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/parser"
)

type CodeRequest struct {
	Code string `json:"code"`
}

type CompileResponse struct {
	Bytecode     []byte        `json:"bytecode"`
	Constants    []interface{} `json:"constants"`
	Instructions string        `json:"instructions"`
	Error        string        `json:"error,omitempty"`
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
