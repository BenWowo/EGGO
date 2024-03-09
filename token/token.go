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

	EQ = "=="
	NE = "!="
	LT = "<"
	GT = ">"
	LE = "<="
	GE = ">="

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

type opPrecPair struct {
	Operator   string
	Precedence int
}

var OpPrecTable = map[string]opPrecPair{
	EQ: {
		Operator:   EQ,
		Precedence: 10,
	},
	NE: {
		Operator:   NE,
		Precedence: 10,
	},
	LT: {
		Operator:   LT,
		Precedence: 11,
	},
	LE: {
		Operator:   LE,
		Precedence: 11,
	},
	GT: {
		Operator:   GT,
		Precedence: 11,
	},
	GE: {
		Operator:   GE,
		Precedence: 11,
	},
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
