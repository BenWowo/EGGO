package scanner

import (
	"eggo/lib"
	"eggo/token"
	"fmt"
)

type Scanner struct {
	content      string
	ch           byte
	position     int
	readPosition int
}

// Returns a pointer to an instance of a Scanner.
func New(filepath string) *Scanner {
	s := &Scanner{
		// TODO - used buffered file reading vs reading in the whole file
		content: lib.FileToString(filepath),
	}
	s.readChar()
	return s
}

// Debug fucntion that prints all of the tokens in the input file.
func (s *Scanner) ScanFile() {
	for tok := s.NextToken(); tok.Type != token.EOF; tok = s.NextToken() {
		fmt.Printf("%v \n", tok)
	}
}

// Returns the next token from the input stream.
// If the end of the input stream has been reached
// NextToken will return an EOF token.
func (s *Scanner) NextToken() token.Token {
	var tok token.Token

	// make this just advance to the next meaningful token
	// within this function I can skip the whitespace and commnets
	s.advancePosition()

	// TODO - handle non space separated tokens

	switch s.ch {
	case '(':
		tok = token.Token{Type: token.LParen, Literal: string(s.ch)}
	case ')':
		tok = token.Token{Type: token.RParen, Literal: string(s.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(s.ch)}
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: string(s.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(s.ch)}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(s.ch)}
	case '*':
		tok = token.Token{Type: token.STAR, Literal: string(s.ch)}
	case '!':
		if s.peekChar() == '=' {
			s.readChar()
			tok = token.Token{Type: token.NE, Literal: token.NE}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(s.ch)}
		}
		// bitwise not
	case '=':
		if s.peekChar() == '=' {
			s.readChar()
			tok = token.Token{Type: token.EQ, Literal: token.EQ}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(s.ch)}
		}
	case '<':
		if s.peekChar() == '<' {
			s.readChar()
			tok = token.Token{Type: token.LSHIFT, Literal: token.LSHIFT}
		} else if s.peekChar() == '=' {
			s.readChar()
			tok = token.Token{Type: token.LE, Literal: token.LE}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(s.ch)}
		}
	case '>':
		if s.peekChar() == '>' {
			s.readChar()
			tok = token.Token{Type: token.RSHIFT, Literal: token.RSHIFT}
		} else if s.peekChar() == '=' {
			s.readChar()
			tok = token.Token{Type: token.GE, Literal: token.GE}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(s.ch)}
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isNumber(s.ch) {
			numberLiteral := s.readNumber()
			tok = token.Token{Type: token.NUMBER_INT, Literal: numberLiteral}
		} else if isLetter(s.ch) {
			ident := s.readIdent()
			if keyword := token.KeywordTable[ident]; keyword != "" {
				tok = token.Token{Type: token.TokenType(keyword), Literal: string(keyword)}
			} else {
				tok = token.Token{Type: token.IDENT, Literal: ident}
			}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(s.ch)}
		}
	}

	s.readChar()
	return tok
}

// Returns the next chracter from the input stream.
// If the end of the input stream has been reached
// readChar will return 0.
func (s *Scanner) readChar() {
	if s.readPosition >= len(s.content) {
		s.ch = 0
	} else {
		// hmm how will this work on a buffered read
		// will I need to call another buffer read when I run out
		// of characters in the buffer
		s.ch = s.content[s.readPosition]
	}

	s.position = s.readPosition
	s.readPosition += 1
}

func (s *Scanner) peekChar() byte {
	if s.readPosition >= len(s.content) {
		return 0
	} else {
		return s.content[s.readPosition]
	}
}

func (s *Scanner) advancePosition() {
	for {
		if s.ch == '/' && s.peekChar() == '/' {
			// TODO add newline escape using backslash '\'
			s.skipComment()
		} else if isWhiteSpace(s.ch) {
			s.skipWhitespace()
		} else {
			break
		}
	}
}

func (s *Scanner) skipComment() {
	for s.ch != '\n' && s.ch != 0 {
		s.readChar()
	}
}

func (s *Scanner) skipWhitespace() {
	for isWhiteSpace(s.ch) {
		s.readChar()
	}
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (s *Scanner) readIdent() string {
	left := s.position
	for isLetter(s.ch) || isNumber(s.ch) {
		s.readChar()
	}
	return s.content[left:s.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (s *Scanner) readNumber() string {
	left := s.position
	for isNumber(s.ch) {
		s.readChar()
	}
	return s.content[left:s.position]
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
