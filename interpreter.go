package main

import (
	"eggo/ast"
	"eggo/token"
	"strconv"
)

func interpretAST(root *ast.ASTnode) float64 {
	if root.IsTerminal {
		value, err := strconv.ParseFloat(root.Token.Literal, 64)
		if err != nil {
			panic(err)
		}
		return value
	}

	switch root.Token.Type {
	case token.PLUS:
		left, right := 0.0, 0.0
		if root.Left != nil {
			left = interpretAST(root.Left)
		}
		if root.Right != nil {
			right = interpretAST(root.Right)
		}
		return left + right
	case token.MINUS:
		left, right := 0.0, 0.0
		if root.Left != nil {
			left = interpretAST(root.Left)
		}
		if root.Right != nil {
			right = interpretAST(root.Right)
		}
		return left - right
	case token.STAR:
		left, right := 1.0, 1.0
		if root.Left != nil {
			left = interpretAST(root.Left)
		}
		if root.Right != nil {
			right = interpretAST(root.Right)
		}
		return left * right
	case token.SLASH:
		left, right := 1.0, 1.0
		if root.Left != nil {
			left = interpretAST(root.Left)
		}
		if root.Right != nil {
			right = interpretAST(root.Right)
		}
		return left / right
	case token.LSHIFT:
		left, right := 0, 0
		if root.Left != nil {
			left = int(interpretAST(root.Left))
		}
		if root.Right != nil {
			right = int(interpretAST(root.Right))
		}
		return float64(left << right)
	case token.RSHIFT:
		left, right := 0, 0
		if root.Left != nil {
			left = int(interpretAST(root.Left))
		}
		if root.Right != nil {
			right = int(interpretAST(root.Right))
		}
		return float64(left >> right)
	default:
		panic("This isn't supposed to happen\n")
	}
}
