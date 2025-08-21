package lexer

import (
	"testing"
	"monkey-playground-backend/token"
)

func TestNextToken_SimpleSequence(t *testing.T) {
	input := `let five = 5; let ten = 10;`

	l := New(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_OperatorsAndDelimiters(t *testing.T) {
	input := `!-/*5; 5 < 10 > 5; 1 == 1; 2 != 3; { } ( ) , ;`

	l := New(input)

	types := []token.TokenType{
		token.BANG, token.MINUS, token.SLASH, token.ASTERISK, token.INT, token.SEMICOLON,
		token.INT, token.LT, token.INT, token.GT, token.INT, token.SEMICOLON,
		token.INT, token.EQ, token.INT, token.SEMICOLON,
		token.INT, token.NOT_EQ, token.INT, token.SEMICOLON,
		token.LBRACE, token.RBRACE, token.LPAREN, token.RPAREN, token.COMMA, token.SEMICOLON,
		token.EOF,
	}

	for i, expected := range types {
		tok := l.NextToken()
		if tok.Type != expected {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)", i, expected, tok.Type, tok.Literal)
		}
	}
}

func TestNextToken_KeywordsAndIdentifiers(t *testing.T) {
	input := `let add = fn(x, y) { return x + y; }; let result = add(5, 10); true != false; if (5 < 10) { return true; } else { return false; }`

	l := New(input)

	expectedTypes := []token.TokenType{
		token.LET, token.IDENT, token.ASSIGN, token.FUNCTION, token.LPAREN, token.IDENT, token.COMMA, token.IDENT, token.RPAREN,
		token.LBRACE, token.RETURN, token.IDENT, token.PLUS, token.IDENT, token.SEMICOLON, token.RBRACE, token.SEMICOLON,
		token.LET, token.IDENT, token.ASSIGN, token.IDENT, token.LPAREN, token.INT, token.COMMA, token.INT, token.RPAREN, token.SEMICOLON,
		token.TRUE, token.NOT_EQ, token.FALSE, token.SEMICOLON,
		token.IF, token.LPAREN, token.INT, token.LT, token.INT, token.RPAREN, token.LBRACE, token.RETURN, token.TRUE, token.SEMICOLON, token.RBRACE,
		token.ELSE, token.LBRACE, token.RETURN, token.FALSE, token.SEMICOLON, token.RBRACE,
		token.EOF,
	}

	for i, expected := range expectedTypes {
		tok := l.NextToken()
		if tok.Type != expected {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)", i, expected, tok.Type, tok.Literal)
		}
	}
}

func TestNextToken_Comments(t *testing.T) {
	input := `// This is a comment
let five = 5; // Another comment
// Full line comment
let ten = 10;
// Final comment`

	l := New(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_CommentsWithDivision(t *testing.T) {
	input := `5 / 2 // Division followed by comment
// Comment line
let x = 10 / 5;`

	l := New(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "5"},
		{token.SLASH, "/"},
		{token.INT, "2"},
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

