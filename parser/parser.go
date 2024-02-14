package parser

import (
	"eggo/scanner"
	"eggo/token"
	"fmt"
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

func (p *Parser) ParseStatement() *ASTnode {
	node := new(ASTnode)

	if p.peekToken.Type == token.EOF {
		node = nil
	} else if p.peekToken.Type == token.PRINT {
		node = p.parsePrintStatement()
	} else {
		fmt.Printf("Unexpcted token Type in parser %s\n", p.peekToken.Type)
		panic("err")
	}

	return node
}

// printStatement: "print" expression ";"
// TODO - create new ast node with multiple children for variadic print args
func (p *Parser) parsePrintStatement() *ASTnode {
	p.nextToken()

	node := &ASTnode{
		Token: p.curToken,
		// Left:  p.ParseBinaryOperation(0), // This caused me a huge headache ;(
	}
	node.Left = p.ParseBinaryOperation(0)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		fmt.Printf("next token: %s\n", p.peekToken.Type)
	} else {
		fmt.Printf("expected semicolon after \"print\" \"expression\" got: %s\n", p.peekToken.Type)
		panic("err")
	}

	return node
}

func (p *Parser) ParseBinaryOperation(previous_precedence int) *ASTnode {
	node := new(ASTnode)

	node.Left = p.parseTerminalNode()

	// TODO handle end of file. As it is now the code will hang
	if p.peekToken.Type == token.SEMICOLON {
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

		if p.peekToken.Type == token.SEMICOLON {
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
