package ast

import (
	"fmt"
	"strings"
)

type ASTnode interface {
	// String() string
}

type DeclareNode struct {
	Ident    string
	DataType string
}

type AssignNode struct {
	Ident      string
	Expression *ExpressionNode
}

type PrintNode struct {
	Expression *ExpressionNode
}

type ExpressionNode struct {
	Value string
	Left  *ExpressionNode
	Right *ExpressionNode
}

func (node *ExpressionNode) IsTerminal() bool {
	return node.Left == nil && node.Right == nil
}

func (node *ExpressionNode) String() string {
	var formatedString func(node *ExpressionNode, level int) string
	formatedString = func(node *ExpressionNode, level int) string {
		if node == nil {
			return "nil"
		}

		prev_indent := strings.Repeat("  ", level)
		cur_indent := strings.Repeat("  ", level+1)

		lBrace := "{\n"
		value := fmt.Sprintf("%sToken: %s\n", cur_indent, node.Value)
		left := fmt.Sprintf("%sLeft: %s\n", cur_indent, formatedString(node.Left, level+1))
		right := fmt.Sprintf("%sRight: %s\n", cur_indent, formatedString(node.Right, level+1))
		rBrace := fmt.Sprintf("%s}", prev_indent)

		return lBrace + value + left + right + rBrace
	}
	return formatedString(node, 0)
}

type BlockNode struct {
	Statements []*ASTnode
}

type IfNode struct {
	Condition *ExpressionNode
	HappyBody *BlockNode
	SadBody   *ASTnode
}

func (node *IfNode) ContainsElse() bool {
	return node.SadBody != nil
}

type WhileNode struct {
	Condition *ExpressionNode
	Body      *BlockNode
}
