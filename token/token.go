package token

type TokenType string

const (
	// operators
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"
	LSHIFT = "<<"
	RSHIFT = ">>"

	LT    = "<"
	GT    = ">"
	LT_EQ = "<="
	GT_EQ = ">="

	// ?
	INT   = "INT"
	IDENT = "IDENT"

	// ?
	SEMICOLON = ";"

	// keywords
	PRINT = "print"

	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var KeywordTable = map[string]TokenType{
	PRINT: PRINT,
}

type Token struct {
	Type    TokenType
	Literal string
}

func Precedence_lookup(tok Token) int {
	switch tok.Type {
	case PLUS:
		return 12
	case MINUS:
		return 12
	case STAR:
		return 13
	case SLASH:
		return 13
	case LSHIFT:
		return 14
	case RSHIFT:
		return 14
	default:
		return 0
	}
}
