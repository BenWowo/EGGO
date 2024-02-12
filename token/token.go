package token

type TokenType string

const (
	PLUS  = "+"
	MINUS = "-"
	STAR  = "*"
	SLASH = "/"

	INT   = "INT"
	IDENT = "IDENT"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

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
	default:
		return 0
	}
}
