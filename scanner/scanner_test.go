package scanner

import (
	"eggo/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "../lib/inputfile.txt"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// declarations
		{token.INT, "int"},
		{token.IDENT, "fred"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		// assignments
		{token.IDENT, "fred"},
		{token.ASSIGN, "="},
		{token.NUMBER_INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "jim"},
		{token.ASSIGN, "="},
		{token.NUMBER_INT, "7"},
		{token.PLUS, "+"},
		{token.IDENT, "fred"},
		{token.SEMICOLON, ";"},

		// print expression
		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.PLUS, "+"},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		// conditionals
		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.EQ, "=="},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.NE, "!="},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.LT, "<"},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.LE, "<="},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.GT, ">"},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},

		{token.PRINT, "print"},
		{token.IDENT, "fred"},
		{token.GE, ">="},
		{token.IDENT, "jim"},
		{token.SEMICOLON, ";"},
	}

	s := New(input)
	for i, tt := range tests {
		tok := s.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tokenLiteral wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestScanFile(t *testing.T) {
	input := "../lib/inputfile.txt"
	s := New(input)
	s.ScanFile()
}
