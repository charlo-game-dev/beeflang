package token

// TokenType represents the type of a token
type TokenType string

// Token represents a lexical token with position information
type Token struct {
	Type    TokenType
	Literal string
	Line    int // line number in source (for error reporting)
	Column  int // column number in source (for error reporting)
}

// Token types
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"  // variable names, function names
	INT    TokenType = "INT"    // integer literals
	STRING TokenType = "STRING" // string literals

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"

	// Comparison operators
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LTE    TokenType = "<="
	GTE    TokenType = ">="

	// Logical operators
	AND TokenType = "&&"
	OR  TokenType = "||"
	NOT TokenType = "!"

	// Delimiters
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	COLON  TokenType = ":"
	COMMA  TokenType = ","
	DOT    TokenType = "."

	// Keywords
	PRAISE      TokenType = "PRAISE"      // function declaration
	BEEF        TokenType = "BEEF"        // block terminator
	FEAST_WHILE TokenType = "FEAST_WHILE" // while loop
	IF          TokenType = "IF"
	ELSE        TokenType = "ELSE"
	PREP        TokenType = "PREP"    // variable declaration
	SERVE       TokenType = "SERVE"   // return
	WRANGLE     TokenType = "WRANGLE" // import module
	HERD        TokenType = "HERD"    // module keyword
	TRUE        TokenType = "TRUE"
	FALSE       TokenType = "FALSE"
	AND_WORD    TokenType = "AND" // 'and' keyword
	OR_WORD     TokenType = "OR"  // 'or' keyword
	NOT_WORD    TokenType = "NOT" // 'not' keyword
)

var keywords = map[string]TokenType{
	"praise":  PRAISE,
	"beef":    BEEF,
	"feast":   FEAST_WHILE, // Will need special handling for "feast while"
	"while":   FEAST_WHILE,
	"if":      IF,
	"else":    ELSE,
	"prep":    PREP,
	"serve":   SERVE,
	"wrangle": WRANGLE,
	"herd":    HERD,
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
