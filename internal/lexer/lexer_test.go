package lexer

import (
	"testing"

	"github.com/elitwilson/beeflang/internal/token"
	"github.com/stretchr/testify/assert"
)

// ========================================
// Literals
// ========================================

func TestTokenizeIntegers(t *testing.T) {
	input := "42"
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.INT, tok.Type)
	assert.Equal(t, "42", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

func TestTokenizeStringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"Hello, Beef!"`, "Hello, Beef!"},
		{`""`, ""},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		assert.Equal(t, token.STRING, tok.Type, "Input: %s", tt.input)
		assert.Equal(t, tt.expected, tok.Literal, "Input: %s", tt.input)
	}
}

// ========================================
// Identifiers
// ========================================

func TestTokenizeIdentifiers(t *testing.T) {
	input := "foo bar"
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.IDENT, tok.Type)
	assert.Equal(t, "foo", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.IDENT, tok.Type)
	assert.Equal(t, "bar", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

func TestTokenizeIdentifiersWithDigits(t *testing.T) {
	// Identifiers can contain digits (but not start with them)
	tests := []struct {
		input    string
		expected string
	}{
		{"result1", "result1"},
		{"var2", "var2"},
		{"test123", "test123"},
		{"a1b2c3", "a1b2c3"},
		{"sum_to_n", "sum_to_n"},
		{"value_42", "value_42"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		assert.Equal(t, token.IDENT, tok.Type, "Input: %s should be IDENT", tt.input)
		assert.Equal(t, tt.expected, tok.Literal, "Input: %s literal mismatch", tt.input)

		// Should be followed by EOF, not another token
		tok = l.NextToken()
		assert.Equal(t, token.EOF, tok.Type, "Input: %s should be followed by EOF", tt.input)
	}
}

// ========================================
// Keywords
// ========================================

func TestTokenizeKeywords(t *testing.T) {
	input := "prep praise beef serve if else"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PREP, "prep"},
		{token.PRAISE, "praise"},
		{token.BEEF, "beef"},
		{token.SERVE, "serve"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

func TestHandleFeastWhileKeyword(t *testing.T) {
	input := "feast while"
	l := New(input)

	// For now, we'll treat "feast while" as two separate tokens
	// This test documents the behavior - we can refine later if needed
	tok := l.NextToken()
	assert.Equal(t, token.FEAST_WHILE, tok.Type)
	assert.Equal(t, "feast", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.FEAST_WHILE, tok.Type)
	assert.Equal(t, "while", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

func TestTokenizeWrangleKeyword(t *testing.T) {
	input := "wrangle"
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.WRANGLE, tok.Type)
	assert.Equal(t, "wrangle", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

func TestTokenizeHerdKeyword(t *testing.T) {
	input := "herd"
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.HERD, tok.Type)
	assert.Equal(t, "herd", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

// ========================================
// Operators
// ========================================

func TestTokenizeSingleCharOperators(t *testing.T) {
	input := "+ - * / = < >"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.ASSIGN, "="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

func TestTokenizeTwoCharOperators(t *testing.T) {
	input := "== != <= >= && ||"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.EQ, "=="},
		{token.NOT_EQ, "!="},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

func TestTokenizeDotOperator(t *testing.T) {
	input := "."
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.DOT, tok.Type)
	assert.Equal(t, ".", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

// ========================================
// Delimiters
// ========================================

func TestTokenizeDelimiters(t *testing.T) {
	input := "( ) : ,"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COLON, ":"},
		{token.COMMA, ","},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

// ========================================
// Whitespace and Comments
// ========================================

func TestSkipWhitespace(t *testing.T) {
	input := "   42   "
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.INT, tok.Type)
	assert.Equal(t, "42", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

func TestSkipComments(t *testing.T) {
	input := `# This is a comment
42  # inline comment`
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.INT, tok.Type)
	assert.Equal(t, "42", tok.Literal)

	tok = l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}

// ========================================
// Integration Tests
// ========================================

func TestSimpleVariableDeclaration(t *testing.T) {
	// Integration test: complete statement
	input := "prep x = 42"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PREP, "prep"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "42"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

func TestTokenizeModuleImportStatement(t *testing.T) {
	input := "wrangle io"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.WRANGLE, "wrangle"},
		{token.IDENT, "io"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

func TestTokenizeMemberAccess(t *testing.T) {
	input := "io.preach"
	l := New(input)

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "io"},
		{token.DOT, "."},
		{token.IDENT, "preach"},
		{token.EOF, ""},
	}

	for i, expected := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, expected.expectedType, tok.Type, "token %d type mismatch", i)
		assert.Equal(t, expected.expectedLiteral, tok.Literal, "token %d literal mismatch", i)
	}
}

// ========================================
// Position Tracking
// ========================================

func TestTrackLineAndColumn(t *testing.T) {
	// Critical for error messages
	input := `prep x = 42
prep y = 5`
	l := New(input)

	// First line tokens
	l.NextToken()        // cut
	l.NextToken()        // x
	l.NextToken()        // =
	tok := l.NextToken() // 42
	assert.Equal(t, 1, tok.Line, "first line should be 1")

	// Second line tokens
	tok = l.NextToken() // prep on line 2
	assert.Equal(t, 2, tok.Line, "second line should be 2")
}

func TestEOFToken(t *testing.T) {
	input := ""
	l := New(input)

	tok := l.NextToken()
	assert.Equal(t, token.EOF, tok.Type)
}
