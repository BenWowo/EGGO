package parser

import (
	"eggo/ast"
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

func (p *Parser) ParseStatement() *ast.ASTnode {
	node := new(ast.ASTnode)

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
func (p *Parser) parsePrintStatement() *ast.ASTnode {
	p.nextToken()

	node := &ast.ASTnode{
		Token: p.curToken,
		// Left:  p.ParseBinaryOperation(0), // This caused me a huge headache ;(
	}
	node.Left = p.ParseBinaryOperation(0)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	} else {
		fmt.Printf("expected semicolon after \"print\" \"expression\" got: %s\n", p.peekToken.Type)
		panic("err")
	}

	return node
}

func (p *Parser) ParseBinaryOperation(previous_precedence int) *ast.ASTnode {
	node := new(ast.ASTnode)

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

		node = &ast.ASTnode{
			Token: p.curToken,
			Left:  prev,
		}

		if p.peekToken.Type == token.SEMICOLON {
			return node.Left
		}
	}

	return node.Left
}

func (p *Parser) parseTerminalNode() *ast.ASTnode {
	p.nextToken()
	return &ast.ASTnode{Token: p.curToken, IsTerminal: true}
}

func (p *Parser) parseOperator() *ast.ASTnode {
	p.nextToken()
	return &ast.ASTnode{Token: p.curToken}
}
