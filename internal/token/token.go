package token

// TokenType represents the type of a token
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
}

// Token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // variable names, function names
	INT   = "INT"   // integer literals

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"

	// Comparison operators
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LTE    = "<="
	GTE    = ">="

	// Logical operators
	AND = "&&"
	OR  = "||"
	NOT = "!"

	// Delimiters
	LPAREN = "("
	RPAREN = ")"
	COLON  = ":"
	COMMA  = ","

	// Keywords
	PRAISE      = "PRAISE"      // function declaration
	BEEF        = "BEEF"        // block terminator
	FEAST_WHILE = "FEAST_WHILE" // while loop
	IF          = "IF"
	ELSE        = "ELSE"
	CUT         = "CUT"    // variable declaration
	SERVE       = "SERVE"  // return
	GENESIS     = "GENESIS" // main/entry point
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	AND_WORD    = "AND" // 'and' keyword
	OR_WORD     = "OR"  // 'or' keyword
	NOT_WORD    = "NOT" // 'not' keyword
)

var keywords = map[string]TokenType{
	"praise":  PRAISE,
	"beef":    BEEF,
	"feast":   FEAST_WHILE, // Will need special handling for "feast while"
	"while":   FEAST_WHILE,
	"if":      IF,
	"else":    ELSE,
	"cut":     CUT,
	"serve":   SERVE,
	"genesis": GENESIS,
	"true":    TRUE,
	"false":   FALSE,
	"and":     AND_WORD,
	"or":      OR_WORD,
	"not":     NOT_WORD,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
