package scanner

import (
	"eggo/token"
	"fmt"
	"log"
	"os"
)

// TODO - used buffered file reading vs reading in the whole file
func fileToString(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %s\n", err)
	}

	contentString := string(content)

	return contentString
}

type Scanner struct {
	content      string
	ch           byte
	position     int
	readPosition int
}

func New(filepath string) *Scanner {
	s := &Scanner{
		content: fileToString(filepath),
	}
	s.readChar()
	return s
}

func (s *Scanner) ScanFile() {
	for tok := s.NextToken(); tok.Type != token.EOF; tok = s.NextToken() {
		fmt.Printf("%v \n", tok)
	}
}

func (s *Scanner) NextToken() token.Token {
	var tok token.Token

	s.skipWhitespace()

	switch s.ch {
	case '-':
		tok = token.Token{Type: token.MINUS, Literal: string(s.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(s.ch)}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(s.ch)}
	case '*':
		tok = token.Token{Type: token.STAR, Literal: string(s.ch)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isNumber(s.ch) {
			numberLiteral := s.readNumber()
			tok = token.Token{Type: token.INT, Literal: numberLiteral}
		} else if isLetter(s.ch) {
			identLiteral := s.readIdent()
			// TODO - add case for keywords
			tok = token.Token{Type: token.IDENT, Literal: identLiteral}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(s.ch)}
		}
	}

	s.readChar()
	return tok
}

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

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.readChar()
	}
}

func (s *Scanner) readIdent() string {
	left := s.position
	for isLetter(s.ch) {
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
