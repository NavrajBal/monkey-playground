package lexer

import "monkey-playground-backend/token"

type Lexer struct {
	input        string // whole input source
	position     int    // current position in input (points to current char)
	readPosition int    // next reading position (after current char)
	ch           byte   // current char under examination
}

// New constructs a new Lexer for the given input string
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns the next token from the input stream
// It advances the lexer as needed and skips whitespace
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace advances the input past spaces, tabs, newlines, and carriage returns
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment advances the input past a line comment (from // to end of line)
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readChar reads the next character, advancing position and readPosition
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar returns the next byte without advancing the lexer
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) { return 0 }
	return l.input[l.readPosition]
}

// readIdentifier consumes an identifier [a-zA-Z_][a-zA-Z0-9_]* and returns its literal
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { l.readChar() }
	return l.input[position:l.position]
}

// readNumber consumes a contiguous sequence of digits and returns its literal
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) { l.readChar() }
	return l.input[position:l.position]
}

// readString reads until the closing double quote or EOF and returns the substring
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 { break }
	}
	return l.input[position:l.position]
}

// isLetter reports whether ch is a letter or underscore
func isLetter(ch byte) bool { return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' }

// isDigit reports whether ch is an ASCII digit
func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }

// newToken constructs a token from a single-character literal
func newToken(tokenType token.TokenType, ch byte) token.Token { return token.Token{Type: tokenType, Literal: string(ch)} }
