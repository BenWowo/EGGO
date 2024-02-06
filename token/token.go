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
