package parser

import (
	"eggo/scanner"
	"eggo/token"
)

type Parser struct {
	s         *scanner.Scanner
	curToken  token.Token
	peekToken token.Token
}

func New(filepath string) *Parser {
	p := &Parser{
		s: scanner.New(filepath),
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.s.NextToken()
}

// BNF grammer
// expression: INTEGER_LITERAL
// | expression * expression
// | expression / expression
// | expression + expression
// | expression - expression

func (p *Parser) ParseBinaryExpression() *ASTnode {
	node := &ASTnode{}

	// should be a terminal node i.e. "3"
	node.Left = &ASTnode{
		Data:       p.curToken,
		IsTerminal: true,
	}
	p.nextToken()

	// should be an operator i.e. "+"
	node.Data = p.curToken
	p.nextToken()

	if p.peekToken.Type != token.EOF {
		node.Right = p.ParseBinaryExpression()
	}

	return node
}
