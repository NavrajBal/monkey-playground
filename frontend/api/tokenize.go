package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NavrajBal/monkey-lang/lexer"
	"github.com/NavrajBal/monkey-lang/token"
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
