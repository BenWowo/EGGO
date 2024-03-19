package parser

import (
	"eggo/ast"
	"eggo/scanner"
	"eggo/token"
	"fmt"
	"log"
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

func (p *Parser) expectPeek(expectedTypes ...token.TokenType) {
	valid := false
	for _, tokType := range expectedTypes {
		if p.peekToken.Type == tokType {
			valid = true
		}
	}

	// I also wanna know which function call make the peek mad
	if !valid {
		errorStr := fmt.Sprintf("Unexpected peek token type: %s\n", p.peekToken.Type)
		errorStr += fmt.Sprintf("Expected one of the following types [%v]", expectedTypes)
		log.Fatalf("%s\n", errorStr)

	}
}

func (p *Parser) expectParseStatement() {
	p.expectPeek(token.PRINT, token.INT, token.IDENT, token.IF, token.WHILE, token.LBrace)
}

func (p *Parser) ParseStatement() *ast.ASTnode {
	node := new(ast.ASTnode)

	p.expectParseStatement()
	switch p.peekToken.Type {
	case token.EOF:
		return nil
	case token.INT:
		*node = p.parseDeclarationStatement()
	case token.IDENT:
		*node = p.parseAssignmentStatement()
	case token.PRINT:
		*node = p.parsePrintStatement()
	case token.LBrace:
		*node = p.parseBlockStatement()
	case token.IF:
		*node = p.parseIfStatement()
	case token.WHILE:
		*node = p.parseWhileStatement()
	default:
		log.Fatalf("Unexpected token Type in parser %s\n", p.peekToken.Type)
	}

	return node
}

func (p *Parser) expectParseDeclarationStatement() {
	p.expectPeek(token.INT) // expect peek <data type>
}

func (p *Parser) parseDeclarationStatement() *ast.DeclareNode {
	node := new(ast.DeclareNode)

	p.expectParseDeclarationStatement()
	p.nextToken()
	node.DataType = p.curToken.Literal

	p.expectPeek(token.IDENT)
	p.nextToken()
	node.Ident = p.curToken.Literal

	p.expectPeek(token.SEMICOLON)
	p.nextToken()

	return node
}

func (p *Parser) expectParseAssignmentStatement() {
	p.expectPeek(token.IDENT)
}

func (p *Parser) parseAssignmentStatement() *ast.AssignNode {
	node := new(ast.AssignNode)

	p.expectParseAssignmentStatement()
	p.nextToken()
	node.Ident = p.curToken.Literal

	p.expectPeek(token.ASSIGN)
	p.nextToken()

	p.expectParseExpression()
	node.Expression = p.parseExpression(0)

	p.expectPeek(token.SEMICOLON)
	p.nextToken()

	return node
}

func (p *Parser) expectParsePrintStatement() {
	p.expectPeek(token.PRINT)
}

func (p *Parser) parsePrintStatement() *ast.PrintNode {
	node := new(ast.PrintNode)

	p.expectParsePrintStatement()
	p.nextToken()

	node.Expression = p.parseExpression(0)

	p.expectPeek(token.SEMICOLON)
	p.nextToken()

	return node
}

func (p *Parser) expectParseExpression() {
	p.expectPeek(token.LParen, token.IDENT, token.NUMBER_INT)
}

func (p *Parser) parseExpression(previous_precedence int) *ast.ExpressionNode {
	node := new(ast.ExpressionNode)

	p.expectParseExpression()
	p.nextToken()
	if p.curToken.Literal == token.LParen {
		node.Left = p.parseExpression(0)
		p.expectPeek(token.RParen)
		p.nextToken()
	} else if p.curToken.Type == token.IDENT || p.curToken.Type == token.NUMBER_INT {
		node.Left = &ast.ExpressionNode{
			Value: p.curToken.Literal,
		}
	}

	// expect peek semicolon or operator or RParen
	// hmmm this is not necessaraly an expect peek because after a number
	// you can have operator, rparen, or semicolon
	// I guess in the case that it is not one of those you expect a semicolon
	if p.peekToken.Type == token.SEMICOLON || p.peekToken.Type == token.RParen {
		return node.Left
	}

	p.expectPeek(
		token.PLUS, token.MINUS, token.STAR, token.SLASH, token.LSHIFT, token.RSHIFT,
		token.EQ, token.NE, token.LT, token.LE, token.GT, token.GE,
	)
	p.nextToken()
	node.Value = p.curToken.Literal

	current_precedence := token.OpPrecTable[node.Value].Precedence
	for current_precedence > previous_precedence {
		prev := node
		prev.Right = p.parseExpression(current_precedence)

		node = &ast.ExpressionNode{
			Value: p.curToken.Literal,
			Left:  prev,
		}

		if p.peekToken.Type == token.SEMICOLON || p.peekToken.Type == token.RParen {
			return node.Left
		}
	}

	return node.Left
}

func (p *Parser) expectParseBlockStatement() {
	p.expectPeek(token.LBrace)
}

func (p *Parser) parseBlockStatement() *ast.BlockNode {
	node := new(ast.BlockNode)

	p.expectParseBlockStatement()
	p.nextToken()

	for p.peekToken.Type != token.RBrace {
		p.expectParseStatement()
		append(node.Statements, p.ParseStatement())
	}

	p.expectPeek(token.RBrace)
	p.nextToken()

	return node
}

func (p *Parser) expectParseIfStatement() {
	p.expectPeek(token.IF)
}

func (p *Parser) parseIfStatement() *ast.IfNode {
	node := new(ast.IfNode)

	p.expectParseIfStatement()
	p.nextToken()

	p.expectPeek(token.LParen)
	p.nextToken()

	p.expectPeek(token.LParen, token.IDENT, token.NUMBER_INT)
	node.Condition = p.parseExpression(0)

	p.expectPeek(token.RParen)
	p.nextToken()

	p.expectPeek(token.LBrace)
	node.HappyBody = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.expectPeek(token.ELSE)
		p.nextToken()

		// hmmm I want it to work for ifelse
		// does else take a statement instead of a block statement?

		// I think this should just be a statement and not a block statement
		p.expectParseStatement()
		node.SadBody = p.ParseStatement()
	}

	return node
}

func (p *Parser) expectParseWhileStatement() {
	p.expectPeek(token.WHILE)
}

func (p *Parser) parseWhileStatement() *ast.WhileNode {
	node := new(ast.WhileNode)

	p.expectParseWhileStatement()
	p.nextToken()

	p.expectPeek(token.LParen)
	p.nextToken()

	p.expectPeek(token.LParen, token.IDENT, token.NUMBER_INT)
	node.Condition = p.parseExpression(0)

	p.expectPeek(token.RParen)
	p.nextToken()

	p.expectPeek(token.LBrace)
	node.Body = p.parseBlockStatement()

	return node
}
