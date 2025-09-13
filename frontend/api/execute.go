package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/NavrajBal/monkey-lang/compiler"
	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/parser"
	"github.com/NavrajBal/monkey-lang/vm"
)

type ExecuteResponse struct {
	Result string `json:"result"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
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
