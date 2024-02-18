package token

type TokenType string

const (
	// operators
	ASSIGN = "="

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

	LParen = "("
	RParen = ")"

	// token types
	NUMBER_INT = "INT"
	IDENT      = "IDENT"

	// ?
	SEMICOLON = ";"

	// keywords
	PRINT = "print"
	INT   = "int"

	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var KeywordTable = map[string]string{
	PRINT: PRINT,
	INT:   INT,
}

type OperatorStruct struct {
	Operator   string
	Precedence int
}

var OperatorTable = map[string]OperatorStruct{
	PLUS: {
		Operator:   PLUS,
		Precedence: 12,
	},
	MINUS: {
		Operator:   MINUS,
		Precedence: 12,
	},
	STAR: {
		Operator:   STAR,
		Precedence: 13,
	},
	SLASH: {
		Operator:   SLASH,
		Precedence: 13,
	},
	LSHIFT: {
		Operator:   LSHIFT,
		Precedence: 14,
	},
	RSHIFT: {
		Operator:   RSHIFT,
		Precedence: 14,
	},
}

type Token struct {
	Type    TokenType
	Literal string
}
