package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // an illegal char we do not know how to parse
	EOF     = "EOF"     // end of file. tells the parser it can stop

	// Identifiers + literals
	IDENT = "IDENT" // variable names
	INT   = "INT"   // integers

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    TokenType
	Literal string
}
