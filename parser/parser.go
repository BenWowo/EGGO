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

func (p *Parser) expectPeek(tokTypes []token.TokenType) {
	expectedTypes := map[token.TokenType]bool{}
	for _, tokType := range tokTypes {
		expectedTypes[tokType] = true
	}

	valid := false
	for _, tokType := range tokTypes {
		if expectedTypes[tokType] {
			valid = true
		}
	}

	// I also wanna know which function call make the peek mad
	if !valid {
		errorStr := fmt.Sprintf("Unexpected peek token type: %s\n", p.peekToken.Type)
		errorStr += fmt.Sprintf("Expected one of the following types [%v]", tokTypes)
		log.Fatalf("%s\n", errorStr)

	}
}

func (p *Parser) ParseStatement() ast.ASTnode {
	node := new(ast.ASTnode)

	// expect peek
	p.expectPeek([]token.TokenType{token.PRINT, token.INT, token.IDENT})
	switch p.peekToken.Type {
	case token.EOF:
		return nil
	case token.PRINT:
		*node = p.parsePrintStatement()
	case token.INT:
		*node = p.parseDeclarationStatement()
	case token.IDENT:
		*node = p.parseAssignmentStatement()
	default:
		log.Fatalf("Unexpected token Type in parser %s\n", p.peekToken.Type)
	}

	return *node
}

// int jim;
func (p *Parser) parseDeclarationStatement() *ast.DeclareNode {
	node := new(ast.DeclareNode)

	p.expectPeek([]token.TokenType{token.INT}) // expect peek <data type>
	p.nextToken()
	node.DataType = p.curToken.Literal

	p.expectPeek([]token.TokenType{token.IDENT})
	p.nextToken()
	node.Ident = p.curToken.Literal

	p.expectPeek([]token.TokenType{token.SEMICOLON})
	p.nextToken()

	return node
}

// john = 10;
// jim = 7 + john;
func (p *Parser) parseAssignmentStatement() *ast.AssignNode {
	node := new(ast.AssignNode)

	p.expectPeek([]token.TokenType{token.IDENT})
	p.nextToken()
	node.Ident = p.curToken.Literal

	p.expectPeek([]token.TokenType{token.ASSIGN})
	p.nextToken()

	node.Expression = p.parseExpression(0)

	p.expectPeek([]token.TokenType{token.SEMICOLON})
	p.nextToken()

	return node
}

// "print" expression ";"
func (p *Parser) parsePrintStatement() *ast.PrintNode {
	node := new(ast.PrintNode)

	p.expectPeek([]token.TokenType{token.PRINT})
	p.nextToken()

	node.Expression = p.parseExpression(0)

	p.expectPeek([]token.TokenType{token.SEMICOLON})
	p.nextToken()

	return node
}

// 7 + john
// 2 + 3 * 5
func (p *Parser) parseExpression(previous_precedence int) *ast.ExpressionNode {
	node := new(ast.ExpressionNode)

	// fmt.Printf("Cur token when called: %v\n", p.curToken)
	p.expectPeek([]token.TokenType{token.LParen, token.IDENT, token.NUMBER_INT})
	p.nextToken()
	if p.curToken.Literal == token.LParen {
		node.Left = p.parseExpression(0)
		p.expectPeek([]token.TokenType{token.RParen})
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

	p.expectPeek([]token.TokenType{
		token.PLUS, token.MINUS, token.STAR, token.SLASH, token.LSHIFT, token.RSHIFT,
	})
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
