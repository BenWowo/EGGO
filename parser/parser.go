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
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.s.NextToken()
}

func (p *Parser) ParseBinaryOperation(previous_precedence int) *ASTnode {
	node := new(ASTnode)

	node.Left = p.parseTerminalNode()

	if p.peekToken.Type == token.EOF {
		return node.Left
	}

	node.Token = p.parseOperator().Token

	current_precedence := token.Precedence_lookup(node.Token)
	for current_precedence > previous_precedence {
		prev := node
		prev.Right = p.ParseBinaryOperation(current_precedence)

		node = &ASTnode{
			Token: p.curToken,
			Left:  prev,
		}

		if p.peekToken.Type == token.EOF {
			return node.Left
		}
	}

	return node.Left
}

func (p *Parser) parseTerminalNode() *ASTnode {
	p.nextToken()
	return &ASTnode{Token: p.curToken, IsTerminal: true}
}

func (p *Parser) parseOperator() *ASTnode {
	p.nextToken()
	return &ASTnode{Token: p.curToken}
}
