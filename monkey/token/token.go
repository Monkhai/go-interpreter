package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // an illegal char we do not know how to parse
	EOF     = "EOF"     // end of file. tells the parser it can stop

	// Identifiers + literals
	IDENTIFIER = "IDENT"  // variable names
	INT        = "INT"    // integers
	STRING     = "STRING" // strings

	// Operators
	ASSIGN = "="
	// INCR     = "++"
	PLUS  = "+"
	MINUS = "-"
	// DECR     = "--"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="

	LT = "<"
	GT = ">"

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
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
